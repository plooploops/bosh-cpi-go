#!/bin/bash

executable=./cpi

stemcell_cid=$(cat ./create_stemcell | \
   ${executable} | \
   jq -r ".result")
echo "Stemcell is ${stemcell_cid}"

vm_cid=$(cat ./create_vm | \
   jq ".arguments[1] = \"${stemcell_cid}\"" | \
   ${executable} | \
   jq -r ".result")
echo "VM is ${vm_cid}"

disk_cid=$(cat ./create_disk | \
   ${executable} | \
   jq -r ".result")
echo "Disk is ${disk_cid}"

attach_result=$(cat ./attach_disk | \
   jq ".arguments[0]=\"${vm_cid}\"" | \
   jq ".arguments[1]=\"${disk_cid}\"" | \
   ${executable})
echo "Attaching disk is ${attach_result}"
