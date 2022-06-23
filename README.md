
## New Relic Agents Paketo Buildpack is a Cloud Native Buildpack that contributes and configures the New Relic Java, NodeJs or Python Agent.


**Behavior**

  
This buildpack will participate if all the following conditions are met

<br/>
The `$BP_NEW_RELIC_ENABLED` is set to true (defaults to false)
<br/><br/>
For Python applications a Procfile is required at root of your application during build stage, sample file available /resources/Procfile
<br/><br/>

The buildpack will do the following for Java applications:
<br/>
Contributes the New Relic Java Agent to the newrelic_java layer configures `$JAVA_TOOL_OPTIONS` to use it
<br/>
Copies New Relic Java Agent configuration file in /resources/newrelic.yml to the newrelic_java layer
<br/><br/>
The buildpack will do the following for NodeJs applications:
<br/>
Installs New Relic Nodejs agent using npm
<br/>
Contributes the New Relic NodeJs Agent to the newrelic-nodejs layer configures `NODE_PATH` to include the added modules
<br/>
Copies New Relic NodeJs Agent configuration file in /resources/newrelic.js to the root folder of your application
<br/>
Configures your main application module with `require('newrelic');`
<br/><br/>
The buildpack will do the following for Python applications:
<br/>
Installs New Relic Python agent using pip3
<br/>
Copies New Relic Python Agent configuration file in /resources/newrelic.ini to the root folder of your application
<br/>

**Variables**
<br/>

| Key | Description |
|--|--|
| `BP_NEW_RELIC_ENABLED` | Defaults to false - will not participate - as set in buildpack.toml   |
| `NEW_RELIC_APP_NAME` | Defaults to app_name variable set in newrelic.yml, newrelic.js or newrelic.ini respectively   |
| `NEW_RELIC_LICENSE_KEY`  | Required at build time or runtime     |
| `NEW_RELIC_AGENT_ENABLED`  | Defaults to agent_enabled variable set in newrelic.yml, newrelic.js or newrelic.ini respectively |

<br/>
You can override any setting from a system property or in the newrelic.yml or newrelic.js by setting an environment variable.

The environment variable corresponding to a given setting in the config file is the setting name prefixed by NEW_RELIC with all dots (.) and dashes (-) replaced by underscores (_). 

For this to work as part of the **build stage**, you will need to precede the variables with BPE, i.e. `BPE_NEW_RELIC_LOG_LEVEL`

For this to work during **runtime**, simply prefix the setting name, i.e. `NEW_RELIC_LOG_LEVEL`

Please refer New Relic Java, Nodejs and Python agents documentation for more information

https://docs.newrelic.com/docs/apm/agents/nodejs-agent/installation-configuration/nodejs-agent-configuration/#environment

https://docs.newrelic.com/docs/apm/agents/java-agent/configuration/java-agent-configuration-config-file/#Environment_Variables

https://docs.newrelic.com/docs/apm/agents/python-agent/configuration/python-agent-configuration/

<br/>

**Examples**

pack build CONTAINERNAME -p ./PATHTOPYTHONAPP -b paketo-buildpacks/python -b ./PATHTOLOCALBUILDPACK ... 

pack build CONTAINERNAME -p ./PATHTONODEJSAPP -b paketo-buildpacks/nodejs -b ./PATHTOLOCALBUILDPACK ...

pack build CONTAINERNAME -p ./PATHTOJAVAAPP -b -b paketo-buildpacks/java  -b ./PATHTOLOCALBUILDPACK ... 

 <br/>
--env `BP_NEW_RELIC_ENABLED`=true \ -----------> Required condition to build the buildpack as default is false

--env `BPE_NEW_RELIC_APP_NAME`=xxxxxxxxxx \ -----------> Optional can be set on runtime

--env `BPE_NEW_RELIC_LICENSE_KEY`=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx -----------> Optional can be set on runtime
 <br/>

**Configuration settings precedence**
 

![Settings precedence](https://docs.newrelic.com/static/java-config-cascade-bb36c948f6227353b43c253c234092df.png)

With the Java agent, server-side configuration overrides all other settings.
Environment variables override Java system properties.
Java properties override user configuration settings in your newrelic.yml file.
User settings override the newrelic.yml default settings.
<br/>
Please refer to New Relic Java agent documentation for more information

https://docs.newrelic.com/docs/apm/agents/java-agent/configuration/java-agent-configuration-config-file

<br/>
By default the Java agent configuration file will be located at /layers/newrelic_java/nr-agent-java/newrelic.yml and for NodeJs and Python at your application root folder, this can be overwritten at runtime, with configmap for kubernetes deployments.


Please refer to Kubernetes documentation for more information about configmaps

[https://kubernetes.io/docs/concepts/configuration/configmap/](https://kubernetes.io/docs/concepts/configuration/configmap/)
