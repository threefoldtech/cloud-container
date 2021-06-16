root=$1
index=$2
init=$3

bridge="zos0"
tap="cont-${index}"

if ! ip l show $tap > /dev/null; then
    sudo ip tuntap add dev "$tap" mode tap
    sudo ip l set "$tap" master $bridge
    sudo ip l set "$tap" up
fi

socket="/tmp/root.${index}.socket"
sudo virtiofsd-rs --shared-dir ${root} --socket ${socket}  &

if [ -z "${init}" ]; then
    init='init=/bin/sh'
fi

exec sudo cloud-hypervisor \
    --kernel output/kernel \
    --initramfs output/initramfs-linux.img \
    --console off \
    --serial tty \
    --cpus boot=1 \
    --memory size=300M,shared=on \
    --fs tag=/dev/root,socket=${socket}  \
    --net tap=${tap} \
    --cmdline "console=ttyS0 rootfstype=virtiofs root=/dev/root rw ${init}" \
    --rng
