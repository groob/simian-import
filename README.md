# simian-import

`simian-import` is a go binary to upload pkgsinfo to Simian. This enables autopkg integration, even though the uploadpkg view/endpoint for Simian [was deprecated](https://github.com/google/simian/commit/3615fcdf295d46af0205dd9afe9a02a5e711523a). 

To test the upload to your Simian instance, install the [`gcloud` command line tool](https://cloud.google.com/sdk/) and authenticate with the AppEngine instance as your default project. Currently you need to supply a name for the 'pkg' (or dmg) being uploaded as well as the path to the pkgsinfo. 
In the future this may be extended to also upload to a Google Storage bucket. Sooner than that, documentation will show how to auto-upload to a Cloudfront-fronted S3 bucket, and inject the PackageCompleteURL into the `simian-import`-ed pkginfo. Keep an eye out for the autopkg processor!

This was originally created as a hackathon project for [MacDevOpsYVR](https://www.macdevops.ca/) 2017
