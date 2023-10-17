/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

/*
Package fake provides a fake over_client for over_testing.

A fake over_client is backed by its simple object store indexed by GroupVersionResource.
You can create a fake over_client with optional objects.

	over_client := NewClientBuilder().WithScheme(over_scheme).WithObj(initObjs...).Build()

You can invoke the methods defined in the Client interface.

When in doubt, it's almost always better not to use this package and instead use
over_envtest.Environment with a real over_client and API server.

WARNING: ⚠️ Current Limitations / Known Issues with the fake Client ⚠️
  - This over_client does not have a way to inject specific errors to test handled vs. unhandled errors.
  - There is some support for sub resources which can cause issues with tests if you're trying to update
    e.g. metadata and status in the same over_reconcile.
  - No OpenAPI validation is performed when creating or updating objects.
  - ObjectMeta's `Generation` and `ResourceVersion` don't behave properly, Patch or Update
    operations that rely on these fields will fail, or give false positives.
*/
package fake
