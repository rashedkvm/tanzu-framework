// Copyright 2023 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package imgpkgutil

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestReconciler(t *testing.T) {
	RegisterFailHandler(Fail)
	suiteConfig, _ := GinkgoConfiguration()
	suiteConfig.FailFast = true
	RunSpecs(t, "Fetcher Unit Tests", suiteConfig)
}

const validImagesLockYAML = `
---
apiVersion: imgpkg.carvel.dev/v1alpha1
images:
- annotations:
    kbld.carvel.dev/id: projects-stg.registry.vmware.com/tkg/tkr-vsphere-nonparavirt:v1.24.9_vmware.1-tkg.1-zshippable
    kbld.carvel.dev/origins: |
      - resolved:
          tag: v1.24.9_vmware.1-tkg.1-zshippable
          url: projects-stg.registry.vmware.com/tkg/tkr-vsphere-nonparavirt:v1.24.9_vmware.1-tkg.1-zshippable
  image: 10.92.174.209:8443/library/tkg/tkr-vsphere-nonparavirt@sha256:b56a4c11a3eef1d3fef51c66b1571f92d45f17cf11de8f89e7706dbbb9b6a287
kind: ImagesLock
`

var _ = Describe("parseImagesLock", func() {
	var (
		imagesLockBytes  []byte
		expectedImageMap map[string]string
	)
	BeforeEach(func() {
		imagesLockBytes = nil
		expectedImageMap = nil
	})
	Context("valid input", func() {
		BeforeEach(func() {
			imagesLockBytes = []byte(validImagesLockYAML)
			expectedImageMap = map[string]string{
				"projects-stg.registry.vmware.com/tkg/tkr-vsphere-nonparavirt:v1.24.9_vmware.1-tkg.1-zshippable": "10.92.174.209:8443/library/tkg/tkr-vsphere-nonparavirt@sha256:b56a4c11a3eef1d3fef51c66b1571f92d45f17cf11de8f89e7706dbbb9b6a287",
			}
		})
		It("should return the expected map of images", func() {
			imageMap, err := ParseImagesLock(imagesLockBytes)
			Expect(imageMap).To(Equal(expectedImageMap))
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("nil input", func() {
		BeforeEach(func() {
			imagesLockBytes = nil
		})
		It("should return nil", func() {
			Expect(ParseImagesLock(imagesLockBytes)).To(BeNil())
		})
	})

	Context("invalid yaml", func() {
		BeforeEach(func() {
			imagesLockBytes = []byte(`%%%%`)
		})
		It("should return nil", func() {
			imageMap, err := ParseImagesLock(imagesLockBytes)
			Expect(imageMap).To(BeNil())
			Expect(err).To(HaveOccurred())
		})
	})
})

var _ = Describe("resolveImages", func() {
	var (
		imageMap   map[string]string
		bundle     map[string][]byte
		wantBundle map[string][]byte
	)
	BeforeEach(func() {
		imageMap = map[string]string{
			"projects-stg.registry.vmware.com/tkg/tkr-vsphere-nonparavirt:v1.24.9_vmware.1-tkg.1-zshippable": "10.92.174.209:8443/library/tkg/tkr-vsphere-nonparavirt@sha256:b56a4c11a3eef1d3fef51c66b1571f92d45f17cf11de8f89e7706dbbb9b6a287",
		}
		bundle = map[string][]byte{
			"path1": []byte("projects-stg.registry.vmware.com/tkg/tkr-vsphere-nonparavirt:v1.24.9_vmware.1-tkg.1-zshippable"),
			"path2": []byte("image2"),
		}
		wantBundle = map[string][]byte{
			"path1": []byte("10.92.174.209:8443/library/tkg/tkr-vsphere-nonparavirt@sha256:b56a4c11a3eef1d3fef51c66b1571f92d45f17cf11de8f89e7706dbbb9b6a287"),
			"path2": []byte("image2"),
		}
	})

	Context("matches", func() {
		It("should replace the original images with target images in the bundle", func() {
			ResolveImages(imageMap, bundle)
			Expect(bundle).To(Equal(wantBundle))
		})
	})

	Context("no matches", func() {
		BeforeEach(func() {
			imageMap = map[string]string{
				"image3": "image3:v1",
			}
			wantBundle = map[string][]byte{
				"path1": []byte("projects-stg.registry.vmware.com/tkg/tkr-vsphere-nonparavirt:v1.24.9_vmware.1-tkg.1-zshippable"),
				"path2": []byte("image2"),
			}
		})
		It("should not change the bundle", func() {
			ResolveImages(imageMap, bundle)
			Expect(bundle).To(Equal(wantBundle))
		})
	})
})
