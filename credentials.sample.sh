#!/bin/bash

REALM=hkraemer-terraform-provider-swp-acc
export SWP_APPLICATION_USER_USERNAME=application-account-xyz
export SWP_APPLICATION_USER_PASSWORD=ilikerandompasswords-from-usermanagement
export SWP_AUTHENTICATOR_URL=https://auth..../auth/realms/$REALM/
export SWP_AIPE_URL=https://$REALM....
