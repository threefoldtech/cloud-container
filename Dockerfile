FROM archlinux as builder
RUN pacman -Sy
RUN pacman -S --noconfirm linux mkinitcpio
COPY mkinitcpio.conf /root/
COPY initcpio /root/initcpio
# override the original initcpio
COPY init /usr/lib/initcpio
WORKDIR /root
RUN KERNELVERSION=$(ls /lib/modules) mkinitcpio -D /usr/lib/initcpio -D initcpio -v -c mkinitcpio.conf -g initramfs-linux.img

FROM scratch
COPY --from=builder /boot/vmlinuz-linux /root/initramfs-linux.img /
