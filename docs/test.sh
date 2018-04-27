#!/bin/bash

executable=./cpi

stemcell_cid=$(cat ./templates/create_stemcell.json | \
   ${executable} | \
   jq -r ".result")
echo "Stemcell is ${stemcell_cid}"

vm_cid=$(cat ./templates/create_vm.json | \
   jq ".arguments[1] = \"${stemcell_cid}\"" | \
   ${executable} | \
   jq -r ".result")
echo "VM is ${vm_cid}"

disk_cid=$(cat ./templates/create_disk.json | \
   ${executable} | \
   jq -r ".result")
echo "Disk is ${disk_cid}"

attach_result=$(cat ./templates/attach_disk.json | \
   jq ".arguments[0]=\"${vm_cid}\"" | \
   jq ".arguments[1]=\"${disk_cid}\"" | \
   ${executable})
echo "Attaching disk is ${attach_result}"
