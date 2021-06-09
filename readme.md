## Introduction
cloud-container is only a builder for a custom initramfs image that allows running containers on `cloud-hypervisor`. The container root is served over `virtiofs`.

The image does the following for this to work:
- [x] Add all required `virtio` modules to support `virtiofs`, `pci`, and `disks`
- [x] Expose environment variables from `/etc/environment` to the container `entrypoint`. A container manager then can simply write down the /etc/environment file in the container root before booting.
- [ ] Configure network interface via cmdline argument passed to the kernel
- [ ] Pre mount attached disks to configured endpoints

### kernel arguments
- `net_ethX=SPEC` argument to configure each interface *TO BE DEFINED*
- `mnt_vdX=/path` auto mount disk to given end point *TODO*

### Building
**pre-requirements**: `docker`

Run
```bash
./build.sh
```

### Testing container
**pre-requirement**: `virtiofsd`, `cloud-hypervisor`
extract a container root (or mount a container flist)

> to edit the environment variables available to your entrypoint you have to edit `<rootfs>/etc/environment`

Run
```bash
./test.sh <rootfs> [entrypoint]
```

default entrypoint in `/bin/bash`
