#!/usr/bin/env python
import argparse
from rc import gcloud


def copy_build_script_to_machine(*, name, local_path):
    machine = gcloud.get(name)
    if machine is None:
        exit(1)

    p = machine.upload(local_path, f'/home/{machine.username}')
    if p.returncode != 0:
        exit(1)


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument('--name', required=True,
                        help='Name of the gcloud machine')
    parser.add_argument('--local_path', required=True,
                        help='Local path of the build script')

    args = parser.parse_args()
    copy_build_script_to_machine(**vars(args))
