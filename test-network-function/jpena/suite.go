package jpena

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/test-network-function/test-network-function/pkg/config"
	"github.com/test-network-function/test-network-function/test-network-function/common"
	"github.com/test-network-function/test-network-function/test-network-function/identifiers"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/test-network-function/test-network-function/pkg/tnf/dependencies"
	"github.com/test-network-function/test-network-function/pkg/tnf/interactive"
	"github.com/test-network-function/test-network-function/pkg/tnf/testcases"
	"github.com/test-network-function/test-network-function/test-network-function/results"

	utils "github.com/test-network-function/test-network-function/pkg/utils"
)

const (
	VersionRegex = `(?m)Red Hat Enterprise Linux release (\d+)\.(\d+).*$`
)

var (
	ReleaseCommand = fmt.Sprintf("%s '^Red Hat' /etc/redhat-release", dependencies.GrepBinaryName)
)

var _ = ginkgo.Describe(common.JpenaTestKey, func() {
	conf, _ := ginkgo.GinkgoConfiguration()
	if testcases.IsInFocus(conf.FocusStrings, common.JpenaTestKey) {
		env := config.GetTestEnvironment()
		ginkgo.BeforeEach(func() {
			env.LoadAndRefresh()
			gomega.Expect(len(env.PodsUnderTest)).ToNot(gomega.Equal(0))
			gomega.Expect(len(env.ContainersUnderTest)).ToNot(gomega.Equal(0))
		})
		ginkgo.ReportAfterEach(results.RecordResult)
		testID := identifiers.XformToGinkgoItIdentifier(identifiers.TestjpenaRedHatRelease)
		ginkgo.It(testID, func() {
			ginkgo.By("should report a proper Red Hat version")
			for _, cut := range env.ContainersUnderTest {
				testRedHatRelease(cut)
			}
		})
	}
})

func testRedHatRelease(cut *config.Container) {
	podName := cut.Oc.GetPodName()
	containerName := cut.Oc.GetPodContainerName()
	context := interactive.NewContext(cut.Oc.GetExpecter(), cut.Oc.GetErrorChannel())

	ginkgo.By(fmt.Sprintf("%s(%s) is being checked for Red Hat release", podName, containerName))

	var commandErr error
	releaseOutput := utils.ExecuteCommand(ReleaseCommand, common.DefaultTimeout, context, func() {
		commandErr = fmt.Errorf("failed to get Red Hat release for container %s", containerName)
	})

	re := regexp.MustCompile(VersionRegex)
	matched := re.FindStringSubmatch(releaseOutput)
	if matched == nil {
		commandErr = fmt.Errorf("container %s(%s) does not run a RHEL image", podName, containerName)
	} else {
		major, _ := strconv.Atoi(matched[1])
		// minor, _ := strconv.Atoi(matched[2])
		// fmt.Printf("JPENA we got version %d.%d\n", major, minor)

		if major < 8 {
			commandErr = fmt.Errorf("RHEL major version is %d, expected >= 8", major)
		}
	}
	gomega.Expect(commandErr).To(gomega.BeNil())
}
