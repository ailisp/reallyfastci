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
                console.log(a)
                return a
            } catch ({ errors }) {
                return { ...state, errors }
            }
        }
    }
}

export default new BuildComponent().mount('my-app');
