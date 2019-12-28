#!/usr/bin/env python
import argparse
from rc import gcloud


def create_machine(*, name, machine_type, disk_size, image_project, image_family, zone, preemptible):
    machine = gcloud.create(
        name=name,
        machine_type=machine_type,
        disk_size=f'{disk_size}G',
        image_project=image_project,
        image_family=image_family,
        zone=zone,
        preemptible=preemptible,
        firewall_allows=[]
    )
    if machine:
        return machine
    else:
        exit(1)


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument('--name', required=True,
                        help='Name of the gcloud machine')
    parser.add_argument('--machine_type', required=True,
                        help='Machine type of the gcloud machine')
    parser.add_argument('--disk_size', required=True,
                        help='Disk size of the gcloud machine in GB')
    parser.add_argument('--image_project', required=True,
                        help='Image project of disk image')
    parser.add_argument('--image_family', required=True,
                        help='Image family of disk image')
    parser.add_argument('--zone', required=True,
                        help='Zone of the gcloud machine')
    parser.add_argument('--preemptible', action='store_true',
                        help='Whether instance is preemptible (1/5 price but can fail at any time)')

    args = parser.parse_args()
    create_machine(**vars(args))
