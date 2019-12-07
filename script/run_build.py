#!/usr/bin/env python
import argparse
from rc import gcloud, run, handle_stream


def run_build(*, name, commit, local_path):
    machine = gcloud.get(name)
    if machine is None:
        exit(1)
    run(['rm', '-rf', f'build/{commit}'])
    run(['mkdir', f'build/{commit}'])

    def stdout_handler(line):
        with open(f'build/{commit}/stdout.log', 'a') as stdo:
            with open(f'build/{commit}/output_combined.log', 'a') as oc:
                stdo.write(line)
                oc.write(line)

    def stderr_handler(line):
        with open(f'build/{commit}/stderr.log', 'a') as stde:
            with open(f'build/{commit}/output_combined.log', 'a') as oc:
                stde.write(line)
                oc.write(line)

    def exit_handler(exitcode):
        with open(f'build/{commit}/exitcode', 'w') as ec:
            with open(f'build/{commit}/output_combined.log', 'a') as oc:
                ec.write(f'{exitcode}\n')
                oc.write(f'Exit Code: {exitcode}\n')

    q, p = machine.run_stream('bash', input=f'''
./{local_path.split('/')[-1]}
''')
    handle_stream(q, stdout_handler=stdout_handler,
                  stderr_handler=stderr_handler, exit_handler=exit_handler)
    exit(p.returncode)


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument('--name', required=True,
                        help='Name of the gcloud machine to run build')
    parser.add_argument('--commit', required=True,
                        help='Commit hash to run build')
    parser.add_argument('--local_path', required=True,
                        help='Local path of build script')
    args = parser.parse_args()
    run_build(**vars(args))
