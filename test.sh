root=$1
index=$2
init=$3

socket="/tmp/root.${index}.socket"
sudo virtiofsd --socket-path=${socket} -o source=${root} &

if [ -z "${init}" ]; then
    init='init=/sbin/zinit "init"'
fi

exec sudo cloud-hypervisor \
    --kernel output/bzImage \
    --initramfs output/initramfs-linux.img \
    --console off \
    --serial tty \
    --cpus boot=1 \
    --memory size=300M,shared=on \
    --fs tag=/dev/root,socket=${socket}  \
    --net tap=cont0 \
    --cmdline "console=ttyS0 rootfstype=virtiofs root=/dev/root rw ip=192.168.123.${index}/24 gw=192.168.123.1 ${init}" \
    --rng
