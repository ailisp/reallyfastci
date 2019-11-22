#!/usr/bin/env python
import argparse
from rc import gcloud

def delete_machine(*, name):
    machine = gcloud.get(name)
    if machine:
        try:
            machine.delete()
        except:
            exit(1)

if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument('--name', required=True,
                        help='Name of the gcloud machine to delete')
    args = parser.parse_args()
    delete_machine(**vars(args))