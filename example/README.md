### Example

After adding `yamlark` to your PATH, execute it with:

`yamlark example.star` or `./example.star` (note the `shabang`)

It will create a new `deployment_new.yaml` where you can verify the differences:

```
--- deployment.yaml	2025-10-30 00:07:08
+++ deployment_new.yaml	2025-11-01 00:04:04
@@ -2,7 +2,7 @@
   enabled: true
   spec:
     minReadySeconds: 0
-    replicas: 1
+    replicas: 7
     strategy:
       rollingUpdate:
         maxSurge: 25%
@@ -24,4 +24,14 @@
         resources: null
         serviceAccount: null
         serviceAccountName: null
+        volumes:
+        - name: obj1-volume
+          persistentVolumeClaim:
+            claimName: obj1-shared-pvc
+        - name: obj2-volume
+          persistentVolumeClaim:
+            claimName: obj2-shared-pvc
+        - name: obj3-volume
+          persistentVolumeClaim:
+            claimName: obj3-shared-pvc
       volumes: null
```
