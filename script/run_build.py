#!/usr/bin/env python
import argparse
from rc import gcloud, run


def run_build(*, name, commit, local_path):
    machine = gcloud.get('machine')
    if machine is None:
        exit(1)
    run(['mkdir', '-p', f'../build/{commit}'])

    stdo = open(f'../build/{commit}/stdout.log', 'w')
    stde = open(f'../build/{commit}/stderr.log', 'w')
    oc = open(f'../build/{commit}/output_combined.log', 'w')
    ec = open(f'../build/{commit}/exitcode', 'w')
    ex = []
    
    def stdout_handler(line):
        stdo.write(line)
        oc.write(line)

    def stderr_handler(line):
        stde.write(line)
        oc.write(line)

    def exit_handler(exitcode):
        ec.write(exitcode)
        ex.append(exitcode)
        oc.write(f'Exit Code: {exitcode}')
    
    handle_stream(machine.run_stream('bash', input=f'''
./{local_path.split('/')[-1]}
''', stdout_handler=stdout_handler, stderr_handler=stderr_handler, exit_handler=exit_handler))
    exit(ex[0])

if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument('--name', required=True,
                        help='Name of the gcloud machine to run build')
    parser.add_argument('--commit', required=True,
                        help='Commit hash to run build')
    parser.add_argument('--local_path', required=True,
                        help='Local path of build script')
    args = parser.parse_args()
    run_build(**args)