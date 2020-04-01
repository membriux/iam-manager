package config

import (
	"context"
	"github.com/golang/mock/gomock"
	"gopkg.in/check.v1"
	"k8s.io/api/core/v1"
	"strings"
	"testing"
)

type PropertiesSuite struct {
	t        *testing.T
	ctx      context.Context
	mockCtrl *gomock.Controller
}

func TestPropertiesSuite(t *testing.T) {
	check.Suite(&PropertiesSuite{t: t})
	check.TestingT(t)
}

func (s *PropertiesSuite) SetUpTest(c *check.C) {
	s.ctx = context.Background()
	s.mockCtrl = gomock.NewController(s.t)
}

func (s *PropertiesSuite) TearDownTest(c *check.C) {
	s.mockCtrl.Finish()
}

// test local properties for local environment
func (s *PropertiesSuite) TestLoadPropertiesLocalEnvSuccess(c *check.C) {
	Props = nil
	err := LoadProperties("LOCAL")
	c.Assert(err, check.IsNil)
	c.Assert(Props, check.NotNil)
	c.Assert(Props.AWSAccountID(), check.Equals, "123456789012")
}

// test failure when env is not local and cm is empty
// should not return nil pointer
func (s *PropertiesSuite) TestLoadPropertiesFailedNoCM(c *check.C) {
	Props = nil
	err := LoadProperties("")
	c.Assert(err, check.NotNil)
	c.Assert(err.Error(), check.Equals, "config map cannot be nil")
}

func (s *PropertiesSuite) TestLoadPropertiesFailedNilCM(c *check.C) {
	Props = nil
	err := LoadProperties("", nil)
	c.Assert(err, check.NotNil)
	c.Assert(err.Error(), check.Equals, "config map cannot be nil")
}

func (s *PropertiesSuite) TestLoadPropertiesSuccess(c *check.C) {
	Props = nil
	cm := &v1.ConfigMap{
		Data: map[string]string{
			"iam.managed.permission.boundary.policy": "iam-manager-permission-boundary",
			"aws.accountId":                          "123456789012",
			"iam.role.max.limit.per.namespace":       "5",
			"aws.region":                             "us-east-2",
			"webhook.enabled":                        "true",
		},
	}
	err := LoadProperties("", cm)
	c.Assert(err, check.IsNil)
	c.Assert(Props.AWSRegion(), check.Equals, "us-east-2")
	c.Assert(Props.MaxRolesAllowed(), check.Equals, 5)
	c.Assert(Props.IsWebHookEnabled(), check.Equals, true)
	c.Assert(Props.AWSAccountID(), check.Equals, "123456789012")
	c.Assert(strings.HasPrefix(Props.ManagedPermissionBoundaryPolicy(), "arn:aws:iam:"), check.Equals, true)
}

func (s *PropertiesSuite) TestLoadPropertiesSuccessWithDefaults(c *check.C) {
	Props = nil
	cm := &v1.ConfigMap{
		Data: map[string]string{
			"iam.managed.permission.boundary.policy": "iam-manager-permission-boundary",
			"aws.accountId":                          "123456789012",
		},
	}
	err := LoadProperties("", cm)
	c.Assert(err, check.IsNil)
	c.Assert(Props.AWSRegion(), check.Equals, "us-west-2")
	c.Assert(Props.MaxRolesAllowed(), check.Equals, 1)
	c.Assert(Props.ControllerDesiredFrequency(), check.Equals, 300)
	c.Assert(Props.IsWebHookEnabled(), check.Equals, false)
	c.Assert(Props.DeriveNameFromNamespace(), check.Equals, false)
	c.Assert(Props.AWSAccountID(), check.Equals, "123456789012")
	c.Assert(strings.HasPrefix(Props.ManagedPermissionBoundaryPolicy(), "arn:aws:iam:"), check.Equals, true)
}

func (s *PropertiesSuite) TestLoadPropertiesSuccessWithCustom(c *check.C) {
	Props = nil
	cm := &v1.ConfigMap{
		Data: map[string]string{
			"iam.managed.permission.boundary.policy": "iam-manager-permission-boundary",
			"aws.accountId":                          "123456789012",
			"iam.role.derive.from.namespace":         "true",
			"controller.desired.frequency":           "30",
			"iam.role.max.limit.per.namespace":       "5",
		},
	}
	err := LoadProperties("", cm)
	c.Assert(err, check.IsNil)
	c.Assert(Props.MaxRolesAllowed(), check.Equals, 5)
	c.Assert(Props.ControllerDesiredFrequency(), check.Equals, 30)
	c.Assert(Props.DeriveNameFromNamespace(), check.Equals, true)
}

func (s *PropertiesSuite) TestGetAllowedPolicyAction(c *check.C) {
	value := Props.AllowedPolicyAction()
	c.Assert(value, check.NotNil)
}

func (s *PropertiesSuite) TestGetRestrictedPolicyResources(c *check.C) {
	value := Props.RestrictedPolicyResources()
	c.Assert(value, check.NotNil)
}

func (s *PropertiesSuite) TestGetRestrictedS3Resources(c *check.C) {
	value := Props.RestrictedS3Resources()
	c.Assert(value, check.NotNil)
}

func (s *PropertiesSuite) TestGetManagedPolicies(c *check.C) {
	value := Props.ManagedPolicies()
	c.Assert(value, check.NotNil)
}

func (s *PropertiesSuite) TestGetAWSAccountID(c *check.C) {
	value := Props.AWSAccountID()
	c.Assert(value, check.NotNil)
}

func (s *PropertiesSuite) TestGetAWSRegion(c *check.C) {
	value := Props.AWSRegion()
	c.Assert(value, check.NotNil)
}

func (s *PropertiesSuite) TestTrustPolicyARNs(c *check.C) {
	value := Props.TrustPolicyARNs()
	c.Assert(value, check.NotNil)
}

func (s *PropertiesSuite) TestGetManagedPermissionBoundaryPolicy(c *check.C) {
	value := Props.ManagedPermissionBoundaryPolicy()
	c.Assert(value, check.NotNil)
}

func (s *PropertiesSuite) TestIsWebhookEnabled(c *check.C) {
	value := Props.IsWebHookEnabled()
	c.Assert(value, check.Equals, false)
}

func (s *PropertiesSuite) TestDeriveNameFromNamespace(c *check.C) {
	value := Props.DeriveNameFromNamespace()
	c.Assert(value, check.Equals, false)
}

func (s *PropertiesSuite) TestControllerDesiredFrequency(c *check.C) {
	value := Props.ControllerDesiredFrequency()
	c.Assert(value, check.Equals, 0)
}