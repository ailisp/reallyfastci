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
        if (commit != "") {
            try {
                let a = { ...state, ...await build.build(commit), commit };
                if (a.exitcode != null) {
                    return a
                } else {
                    this.run('running-build-log', commit)
                    return { ...state, commit, status }
                }
            } catch ({ errors }) {
                return { ...state, commit, errors }
            }
        }
    }


    @on('running-build-log') runningBuildLog = async (state, commit) => {
        const fetchedResource = await fetch(`${defaultBasePath}/api/build/${commit}/output`)
        const reader = await fetchedResource.body.getReader()
        const decoder = new TextDecoder('utf-8');
        const _this = this;

        reader.read().then(function processText({ done, value }) {
            if (done) {
                return;
            }

            _this.run('new-output', decoder.decode(value));
            reader.read().then(processText);
        })
    }

    @on('new-output') newOutput = async (state, newLog) => {
        return { ...state, output_combined: state.output_combined + newLog }
    }
}

export default new BuildComponent().mount('my-app');
