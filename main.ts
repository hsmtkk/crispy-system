import { Construct } from "constructs";
import { App, TerraformStack, CloudBackend, NamedCloudWorkspace } from "cdktf";
import * as google from '@cdktf/provider-google';

//const default_location = 'asia-northeast1';
const project_id = 'crispy-system-366103';

class MyStack extends TerraformStack {
  constructor(scope: Construct, name: string) {
    super(scope, name);

    new google.GoogleProvider(this, 'Google', {
      project: project_id,
    })

    new google.CloudbuildTrigger(this, 'cloud_build_trigger', {
      filename: 'cloudbuild.yaml',
      github: {
        name: 'crispy-system',
        owner: 'hsmtkk',
        push: {
          branch: 'main',
        },
      },
    });
  }
}

const app = new App();
const stack = new MyStack(app, "crispy-system");
new CloudBackend(stack, {
  hostname: "app.terraform.io",
  organization: "hsmtkkdefault",
  workspaces: new NamedCloudWorkspace("crispy-system")
});
app.synth();
