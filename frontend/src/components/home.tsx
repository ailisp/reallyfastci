import app, { Component, on } from 'apprun';
import { index } from '../api';
import { BuildItem } from '../model';

interface State {
    running: BuildItem[],
    finished: BuildItem[]
}

class HomeComponent extends Component {
    state: State = {
        running: [],
        finished: [],
    };

    view = state => {
        return (<div>
            {state.errors && state.errors.map(e => <p>e</p>)}
            <p>Running:</p>
            <ul>{state.running.map(b => <li>
                <a href={`#build/${b}`}>{b}</a>
            </li>)}</ul>
            <p>Finished:</p>
            <ul>{state.finished.map(b => <li>
                <a href={`#build/${b}`}>{b}</a>
            </li>)}</ul>
        </div>)
    };

    getState = async (state) => {
        try {
            return await index.listbuild();
        } catch ({ errors }) {
            return { ...state, errors }
        }
    }

    @on('#') home = async (_state) => await index.listbuild();
}

export default new HomeComponent().mount('my-app');
