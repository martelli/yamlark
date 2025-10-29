### Example

After adding `yamlark` to your PATH, execute it with:

`yamlark bogus.star` or `./bogus.star` (note the `shabang`)

It will create a new `deployment_new.yaml` where you can verify the differences:

```
--- deployment.yaml	2025-10-30 00:07:08
+++ deployment_new.yaml	2025-10-30 00:07:14
@@ -2,7 +2,7 @@
   enabled: true
   spec:
     minReadySeconds: 0
-    replicas: 1
+    replicas: 7
     strategy:
       rollingUpdate:
         maxSurge: 25%
@@ -24,4 +24,8 @@
         resources: null
         serviceAccount: null
         serviceAccountName: null
+        volumes:
+        - name: efs-volume
+          persistentVolumeClaim:
+            claimName: efs-shared-pvc
       volumes: null
```
