package bundledeployment

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rancher/fleet/internal/cmd/controller/mocks"
	fleet "github.com/rancher/fleet/pkg/apis/fleet.cattle.io/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestSetNamespaceLabelsAndAnnotations(t *testing.T) {
	tests := map[string]struct {
		bd          *fleet.BundleDeployment
		ns          corev1.Namespace
		release     string
		expectedNs  corev1.Namespace
		expectedErr error
	}{
		"NamespaceLabels and NamespaceAnnotations are appended": {
			bd: &fleet.BundleDeployment{Spec: fleet.BundleDeploymentSpec{
				Options: fleet.BundleDeploymentOptions{
					NamespaceLabels:      &map[string]string{"optLabel1": "optValue1", "optLabel2": "optValue2"},
					NamespaceAnnotations: &map[string]string{"optAnn1": "optValue1"},
				},
			}},
			ns: corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"name": "test"},
				},
			},
			release: "namespace/foo/bar",
			expectedNs: corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      map[string]string{"name": "test", "optLabel1": "optValue1", "optLabel2": "optValue2"},
					Annotations: map[string]string{"optAnn1": "optValue1"},
				},
			},
			expectedErr: nil,
		},

		"NamespaceLabels and NamespaceAnnotations removes entries that are not in the options, except the name label": {
			bd: &fleet.BundleDeployment{Spec: fleet.BundleDeploymentSpec{
				Options: fleet.BundleDeploymentOptions{
					NamespaceLabels:      &map[string]string{"optLabel": "optValue"},
					NamespaceAnnotations: &map[string]string{},
				},
			}},
			ns: corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      map[string]string{"nsLabel": "nsValue", "name": "test"},
					Annotations: map[string]string{"nsAnn": "nsValue"},
				},
			},
			release: "namespace/foo/bar",
			expectedNs: corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      map[string]string{"optLabel": "optValue", "name": "test"},
					Annotations: map[string]string{},
				},
			},
			expectedErr: nil,
		},

		"NamespaceLabels and NamespaceAnnotations updates existing values": {
			bd: &fleet.BundleDeployment{Spec: fleet.BundleDeploymentSpec{
				Options: fleet.BundleDeploymentOptions{
					NamespaceLabels:      &map[string]string{"bdLabel": "labelUpdated"},
					NamespaceAnnotations: &map[string]string{"bdAnn": "annUpdated"},
				},
			}},
			ns: corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      map[string]string{"bdLabel": "nsValue"},
					Annotations: map[string]string{"bdAnn": "nsValue"},
				},
			},
			release: "namespace/foo/bar",
			expectedNs: corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      map[string]string{"bdLabel": "labelUpdated"},
					Annotations: map[string]string{"bdAnn": "annUpdated"},
				},
			},
			expectedErr: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDynamic := mocks.NewMockInterface(ctrl)
			mockNamespaceableResourceInterface := mocks.NewMockNamespaceableResourceInterface(ctrl)
			u, _ := runtime.DefaultUnstructuredConverter.ToUnstructured(&test.ns)
			mockDynamic.EXPECT().Resource(gomock.Any()).Return(mockNamespaceableResourceInterface).Times(2)
			mockNamespaceableResourceInterface.EXPECT().List(gomock.Any(), metav1.ListOptions{
				LabelSelector: "name=namespace",
			}).Return(&unstructured.UnstructuredList{
				Items: []unstructured.Unstructured{{Object: u}},
			}, nil).Times(1)
			uns, _ := runtime.DefaultUnstructuredConverter.ToUnstructured(&test.expectedNs)
			mockNamespaceableResourceInterface.EXPECT().Update(gomock.Any(), &unstructured.Unstructured{Object: uns}, gomock.Any()).Times(1)

			h := handler{
				dynamic: mockDynamic,
			}
			err := h.setNamespaceLabelsAndAnnotations(test.bd, test.release)

			if err != test.expectedErr {
				t.Errorf("expected error doesn't match: expected %v, got %v", test.expectedErr, err)
			}
		})
	}
}
