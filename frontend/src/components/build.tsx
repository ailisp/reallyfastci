import app, { Component, on } from 'apprun';
import { build } from '../api';
import { runInThisContext } from 'vm';
import { stat } from 'fs';
declare var defaultBasePath;

interface State {
    commit: string,
    output_combined: string,
    status: string,
    exitcode: null | number,
}

class BuildComponent extends Component {
    eventSource: any = null;

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
        if (this.eventSource == null) {
            this.eventSource = new EventSource(`${defaultBasePath}/sse?stream=build-status`);
            this.eventSource.onmessage = (evt) => this.run('build-status-event', evt.data);
        }
        try {
            let a = await build.build(commit);
            if (a.exitcode != null) {
                return { commit, ...a }
            } else {
                this.run('build-status', a.status)
                return { commit, status: a.status }
            }
        } catch ({ errors }) {
            return { ...state, commit, errors }
        }
    }

    @on('running-build-log') runningBuildLog = async (state) => {
        console.log('running-build-log')
        let { commit } = state;
        const fetchedResource = await fetch(`${defaultBasePath}/api/build/${commit}/output`)
        if (fetchedResource.status == 200) {
            const reader = await fetchedResource.body.getReader()
            const decoder = new TextDecoder('utf-8');
            const _this = this;

            reader.read().then(function processText({ done, value }) {
                if (done) {
                    _this.run('finished-status')
                    return;
                }

                _this.run('new-output', decoder.decode(value));
                reader.read().then(processText);
            })
        } else {
            if (state.status != 'Succeed' && state.status != 'Cancelled' && state.status != 'Failed')
                setInterval(() => this.run('running-build-log'), 5000);
        }
    }

    @on('finished-status') finishedStatus = async (state) => {
        let { commit } = state;
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

    @on('build-status-event') buildStatusEvent = async (state, event) => {
        console.log("build status event: " + event)
        let status = JSON.parse(event)
        if (status.commit == state.commit) {
            this.run('build-status', status.status)
        }
    }

    @on('build-status') buildStatus = async (state, status) => {
        if (status == 'Script Copied' && state.status != 'Script Copied') {
            this.run('running-build-log');
        }
        return { ...state, status }
    }
}

export default new BuildComponent().mount('my-app');
