#!/usr/bin/env yamlark

content = yaml.read(path="deployment.yaml")
content['deployment']['spec']['template']['spec']['volumes'] = [{"name": "efs-volume", "persistentVolumeClaim": {"claimName": "efs-shared-pvc"}}]
content['deployment']['spec']['replicas'] = 7
out = yaml.dump(content)
print(out)

