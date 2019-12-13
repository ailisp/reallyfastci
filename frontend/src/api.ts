import app from 'apprun';
declare var webpackenv;
// Conduit API
window['defaultBasePath'] = webpackenv.API_URL;

import { toQueryString, serializeObject, getToken, setToken, get, post, del, put } from './fetch';
export { toQueryString, serializeObject };
import { BuildItem } from './model';

export interface ListBuildResponse {
    running: BuildItem[],
    finished: BuildItem[],
}

export interface BuildResponse {
    status: string,
    output_combined?: string,
    exitcode?: number
}

export interface FinishedStatus {
    status: string,
    exitcode: number
}

export const index = {
    listbuild: () => get<ListBuildResponse>('/api/build')
}

export const build = {
    build: (commit) => get<BuildResponse>(`/api/build/${commit}`),
    finishedStatus: (commit) => get<FinishedStatus>(`/api/build/${commit}/exitcode`)
}
