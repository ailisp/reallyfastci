#!/usr/bin/env python
import argparse
from rc import gcloud


def clone_repo_on_machine(*, name, repo, url, branch, commit):
    machine = gcloud.get(name)
    if machine is None:
        exit(1)

    checkout_commit = f'git checkout {commit}'
    repo = repo.split('/')[-1]
    p = machine.run('bash', input=f'''
git clone {url} --single-branch --branch {branch} {repo}
cd {repo}
{checkout_commit if commit else ""}
''')
    exit(p.returncode)


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument('--name', required=True,
                        help='Name of the gcloud machine')
    parser.add_argument('--url', required=True,
                        help='Repo url to clone')
    parser.add_argument('--repo', required=True,
                        help='Repo name')
    parser.add_argument("--branch", required=True,
                        help='Branch to clone')
    parser.add_argument('--commit', help="commit hash to checkout")

    args = parser.parse_args()
    clone_repo_on_machine(**vars(args))
