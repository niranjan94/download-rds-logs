
### Usage

**Get access key ID and secret access key**
1. Open the IAM console.
1. From the navigation menu, click Users.
1. Select your IAM user name.
1. Click User Actions, and then click Manage Access Keys.
1. Click Create Access Key.
1. Your keys will look something like this:<br>
    Access key ID example: `AKIAIOSFODNN7EXAMPLE`<br>
    Secret access key example: `wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY`
    
1. Click Download Credentials, and store the keys in a secure location.

**Run the tool with your credentials to download**

```bash
AWS_REGION=<your-region> AWS_ACCESS_KEY_ID=<access-key-id> AWS_SECRET_ACCESS_KEY=<secret-access-key> ./download-rds-logs -instance <rds-instance>
```



#### License

```$xslt
Copyright 2018 Niranjan Rajendran

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```