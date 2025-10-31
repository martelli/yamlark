#!/usr/bin/env yamlark


def mangle():
    tom = toml.read("data.toml")
    content = yaml.read(path="deployment.yaml")
    volumes = []
    for entry in tom["main"]["references"]: 
    	e = {"name": "{e}-volume".format(e=entry), "persistentVolumeClaim": {"claimName": "{e}-shared-pvc".format(e=entry)}}
	volumes.append(e)
    content['deployment']['spec']['replicas'] = 7
    content['deployment']['spec']['template']['spec']['volumes'] = volumes
    out = yaml.dumps(content)
    print(out)
    file.write("deployment_new.yaml", out)

mangle()
