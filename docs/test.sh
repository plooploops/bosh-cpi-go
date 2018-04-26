#!/bin/bah

stemcell_cid = cat create_stemcell | ./cpi | jq -r ".result"
vm_cid = cat create_vm | ./cpi | jq -r ".result"
disk_cid = cat create_disk | ./cpi | jq -r ".result"
cat attach_disk | ./cpi
