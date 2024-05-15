# technical-test number 1
## Challenge
The following test will require you to do the following:
- Create a simple application which has a single "/healthcheck" endpoint.
- Containerise your application as a single deployable artifact, encapsulating all dependencies.
- Create a CI pipeline for your application

The application can be written in any programming language. In this solution it's done with GoLang.

The application should be a simple, small, operable web-style API or service provider. It should implement the following:
- An endpoint which returns basic information about your application in JSON format which is generated; The following is expected:
  - Applications Version.
  - Description. ("static variable")
  - Last Commit SHA.

### API Example Response
```
"myapplication": [
  {
    "version": "1.0",
    "description" : "pre-interview technical test",
    "lastcommitsha": "abc57858585"
  }
]
```

The application should have a CI pipeline that is executed when new code is commit and pushed, this pipeline should be comprehensive and cover aspects such as quality, and security; Travis or similar, for example.

Other things to consider as additions:
- Create tests or a test suite; the type of testing is up to you.
- Describe or demonstrate any risks associated with your application/deployment.
- Write a clear and understandable README which explains your application and its deployment steps.

The application code should be within a Github or Gitlab repository where we can review your source code and any configuration required for your project to execute. Please make sure the repository is public so it's viewable.  Below, and in this repository, contains an example solution.

## Example Solution
### Prerequisites

The CI/CD pipeline has been setup with my own GCP account. To set it up with your GCP account, you will need to run following gcloud commands:
```
gcloud auth login
# replace [PROJECT_ID] with your real project ID
gcloud config set project [PROJECT_ID]
gcloud auth configure-docker
gcloud services enable container.googleapis.com cloudbuild.googleapis.com
PROJECT_NUMBER="$(gcloud projects describe ${PROJECT_ID} --format='get(projectNumber)')" \
  gcloud projects add-iam-policy-binding ${PROJECT_NUMBER} \
    --member=serviceAccount:${PROJECT_NUMBER}@cloudbuild.gserviceaccount.com \
    --role=roles/container.developer
```
Then a trigger in CloudBuild can be set up to consume cloudbuild.yaml.

Also a Kubernetes cluster is needed for the CD pipeline to run, for example:
```
gcloud container clusters create cloudbuild-cluster1 \
    --num-nodes=2 --zone=australia-southeast1-c \
    --machine-type=g1-small
```

### Local Build
A bash script is provided to build the app container locally. Unit tests are done within Dockerfile so the build will fail if any unit test didn't pass.
```
./build-local.sh
```

### Local Test
Start the app container:
```
docker run --rm -ti -p 10000:10000 test2
```
Test the endpoint with curl:
```
curl http://localhost:10000/healthcheck
{
    "myapplication": [
        {
            "version": "local-docker",
            "description": "Go test in CloudBuild in multiple steps.",
            "lastcommitsha": "31747bf46587c63040e085b2a854ad9c1a38074d"
        }
    ]
}
```

### Code Linter
I use Golang and Docker linter plugins for Atom so code is checked everytime when it's saved.

### Known Risks
- Fully built on top of GCP toolkit, not portable
- Git branches are not checked for now
- Kubernetes cluster is not checked and assumed in working order within this pipeline
- Commit message is not supported by CloudBuild for now

### CI Pipeline
CI pipeline is built with Google CloudBuild it can be executed within GCP. Looks like the GCB isn't very mature yet as I couldn't find some ENV variables where I can find in other CI tools such as BuildKite.

The workflow of the pipeline is:
- Once a commit has been pushed up, a build in CloudBuild will be kicked off and there are 5 steps in the pipeline:
- TestAndBuild: Unit tests will be run as part of the dockerfile. Then a multi stage docker build will be kicked off. The final artefact will have SHORT_SHA as tag(SHORT_SHA is supplied by GCB's trigger). There's no numeric build number found in GCB so the version currently is fixed to 1.0.
- Tag: also tag this successful build as latest, ie. SHORT_SHA -> latest
- Push: push container images up to Google Container Registry
- DeployApp: create/update the kubernetes deployment and service for the app, this is useful when there's change for the k8s schemas
- RollingUpdate: release the recent built container images to the k8s cluster, in a rolling-update fashion

Then the app can be visited via http://http://35.197.161.61/healthcheck (The IP can be retrieved from GCP K8S console and will change if LoadBalancer has been rebuilt)

### Branching
The trigger in GCB will only be triggered by commits in `dev` and `master` branches. Any other branch won't trigger a build. This is designed to leave room for feature branches and local development. Locally tested feature branch can be merged into `dev` branch via pull requests. `master` branch is considered production so in reality there will be another step to merge `dev` into `master` and also deploy to production cluster.

### Kubernetes
I've included a simple K8s deployment and service to expose this app publicly. By default the latest image will be used:
```
image: asia.gcr.io/idyllic-depth-239301/tech-test2:latest
```
When a build is being deployed by CI, the specific version determined by $SHORT_SHA will be deployed via rolling update:
```
# cloudbuild.yaml
- name: 'gcr.io/cloud-builders/kubectl'
  id: RollingUpdate
  args:
    - 'set'
    - 'image'
    - 'deployment/app-deploy'
    - 'app-golang=asia.gcr.io/$PROJECT_ID/tech-test2:$SHORT_SHA'
```
