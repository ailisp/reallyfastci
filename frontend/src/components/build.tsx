import app, { Component, on } from 'apprun';
import { build } from '../api';

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
                }
            } catch ({ errors }) {
                return { ...state, errors }
            }
        }
    }

    @on('running-build-log') runningBuildLog = async (state, commit) => {

    }
}

export default new BuildComponent().mount('my-app');
