# preview-compare

This tool help visually compare two urls side by side. A common usecase is to deploy preview from a PR and compare new changes with current Prod version. 

# getting started

Note: In example below, both the URL's are `https://rajatjindal.com`, but you get the idea.

```
curl -XPOST 'https://preview-main-wygradga.fermyon.app/api/preview' \
  -H'Content-Type: application/json' \
  -d '{
        "this": "https://rajatjindal.com", 
        "that": "https://rajatjindal.com", 
   }'

{
  "id":"preq-2ec18487-dfb1-4a08-8a9b-2a587ba50945",
  "this":"https://rajatjindal.com",
  "that":"https://rajatjindal.com"
}
```

and now you can go to https://preview-main-wygradga.fermyon.app/api/preview/preq-2ec18487-dfb1-4a08-8a9b-2a587ba50945

on that page, left frame is the 'leader', so any scroll/link clicks you do on left frame, will automatically be reflected in right frame
