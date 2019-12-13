import app, { Component, on } from 'apprun';
import { build } from '../api';
declare var defaultBasePath;

interface State {
    commit: string,
    output_combined: string,
    status: string,
    exitcode: null | number,
}

class BuildComponent extends Component {
    websocket: any;

    state: State = {
        commit: "",
        output_combined: "",
        status: "",
        exitcode: null,
    }

    view = (state) => {
        if (state.commit != "") {
            return (<div>
                <p>Commit: {state.commit}</p>
                <p>Status: {state.status}</p>
                <p>Exit code: {state.exitcode}</p>
                <pre>
                    {state.output_combined}
                </pre>
            </div>)
        }
    }

    @on('#build') build = async (state, commit) => {
        if (state.commit == "" || state.commit != commit) {
            try {
                let a = { ...state, ...await build.build(commit), commit };
                if (a.exitcode != null) {
                    return a
                } else {
                    this.startWebSocket(commit)
                    this.run('running-build-log', commit)
                    return a
                }
            } catch ({ errors }) {
                return { ...state, commit, errors }
            }
        }
    }

    startWebSocket = (commit) => {
        this.websocket = new WebSocket(`${defaultBasePath}/ws`.replace('http', 'ws'));
        this.websocket.onopen = (evt) => {
            this.websocket.send("open")
            console.log("websocket open");
        }
        this.websocket.onclose = (evt) => console.log("websocket close");
        this.websocket.onmessage = (evt) => this.run('ws-msg', evt.data);
        this.websocket.onerror = (evt: MessageEvent) => console.log("websocket error: " + evt.data);
    }


    @on('running-build-log') runningBuildLog = async (state, commit) => {
        const fetchedResource = await fetch(`${defaultBasePath}/api/build/${commit}/output`)
        if (fetchedResource.status == 200) {
            const reader = await fetchedResource.body.getReader()
            const decoder = new TextDecoder('utf-8');
            const _this = this;

            reader.read().then(function processText({ done, value }) {
                if (done) {
                    _this.websocket.close()
                    _this.run('finished-status', commit)
                    return;
                }

                _this.run('new-output', decoder.decode(value));
                reader.read().then(processText);
            })
        }
    }

    @on('finished-status') finishedStatus = async (state, commit) => {
        try {
            let status = await build.finishedStatus(commit);
            return { ...state, status: status.status, exitcode: status.exitcode }
        } catch ({ errors }) {
            return { ...state, errors }
        }
    }

    @on('new-output') newOutput = async (state, newLog) => {
        return { ...state, output_combined: state.output_combined + newLog }
    }

    @on('ws-msg') wsMsg = async (state, message) => {
        console.log("ws message: " + message)
        let status = JSON.parse(message)
        if (status.commit == state.commit) {
            this.run('running-build-log', state.commit)
            return { ...state, status: status.status }
        }
    }
}

export default new BuildComponent().mount('my-app');
