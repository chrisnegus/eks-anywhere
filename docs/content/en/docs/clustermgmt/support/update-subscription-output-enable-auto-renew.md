---
toc_hide: true
---

```
{
    "subscription": {
        "id": "e29fd0d2-d8a8-4ed4-be54-c6c0dd0f7964",
        "arn": "arn:aws:eks:<region>:<account-id>:eks-anywhere-subscription/e29fd0d2-d8a8-4ed4-be54-c6c0dd0f7964",
        "name": "my-subscription",
        "createdAt": "2023-10-10T08:33:36.869000-05:00",
        "effectiveDate": "2023-10-10T08:33:36.869000-05:00",
        "expirationDate": "2024-10-10T08:33:36.869000-05:00",
        "licenseQuantity": 1,
        "licenseType": "CLUSTER",
        "term": {
            "duration": 12,
            "unit": "MONTHS"
        },
        "status": "ACTIVE",
        "packageRegistry": "<packages-registry>.dkr.ecr.<region>.amazonaws.com",
        "autoRenew": true,
        "licenseArns": [
            "arn:aws:license-manager::<account-id>:license:l-4f36acf12e6d491484812927b327c066"
        ],
        "licenses": [
            {
                "id": "<license-id-string>",
                "token": "<license-token-string>"
            }
        ],
        "tags": {
            "environment": "prod"
        }
    }
}
```