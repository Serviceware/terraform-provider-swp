
## 0.2.0

IMPROVEMENTS:

- The OIDC Token is now refreshed whenever it is close to expiry. This ensures apply's
  which run longer than the normal token lifetime (5min usually) are authenticated the
  entire time.

## 0.1.3

FIXES:

- An `swp_data_link_object`, which generates zero links, is now in sync if there are
  no links in the AIPE. Previously this was showing up as a diff due to the difference
  between `null` and and empty set.

## 0.1.0

FEATURES:

- Resource "swp_aipe_data_object" allows you to manage data objects in an AI Process Engine
- Resource "swp_aipe_data_object_link" allows you to manage links between data objects.
- Data Object "swp_aipe_data_object" allows you to fetch information about data objects from
  your AI  Process Instance. 
