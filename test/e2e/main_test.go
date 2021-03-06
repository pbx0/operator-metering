package e2e

import (
	"flag"
	"log"
	"os"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"

	"github.com/operator-framework/operator-metering/pkg/operator"
	"github.com/operator-framework/operator-metering/test/framework"
)

var (
	testFramework *framework.Framework

	reportTestTimeout         = 5 * time.Minute
	reportTestOutputDirectory string
	runAWSBillingTests        bool
)

func init() {
	if reportTestTimeoutStr := os.Getenv("REPORT_TEST_TIMEOUT"); reportTestTimeoutStr != "" {
		var err error
		reportTestTimeout, err = time.ParseDuration(reportTestTimeoutStr)
		if err != nil {
			log.Fatalf("Invalid REPORT_TEST_TIMEOUT: %v", err)
		}
	}
	reportTestOutputDirectory = os.Getenv("TEST_RESULT_REPORT_OUTPUT_DIRECTORY")
	if reportTestOutputDirectory == "" {
		log.Fatalf("$TEST_RESULT_REPORT_OUTPUT_DIRECTORY must be set")
	}

	err := os.MkdirAll(reportTestOutputDirectory, 0777)
	if err != nil {
		log.Fatalf("error making directory %s, err: %s", reportTestOutputDirectory, err)
	}

	runAWSBillingTests = os.Getenv("ENABLE_AWS_BILLING_TESTS") == "true"
}

func TestMain(m *testing.M) {
	kubeconfig := flag.String("kubeconfig", "", "kube config path, e.g. $HOME/.kube/config")
	ns := flag.String("namespace", "metering-ci", "test namespace")
	httpsAPI := flag.Bool("https-api", false, "If true, use https to talk to Metering API")
	flag.Parse()

	var err error
	if testFramework, err = framework.New(*ns, *kubeconfig, *httpsAPI); err != nil {
		logrus.Fatalf("failed to setup framework: %v\n", err)
	}

	os.Exit(m.Run())
}

func TestReportingE2E(t *testing.T) {
	reportsProduceDataTestCases := []reportProducesDataTestCase{
		{
			name:          "namespace-cpu-request",
			queryName:     "namespace-cpu-request",
			newReportFunc: testFramework.NewSimpleReport,
			timeout:       reportTestTimeout,
		},
		{
			name:          "namespace-cpu-usage",
			queryName:     "namespace-cpu-usage",
			newReportFunc: testFramework.NewSimpleReport,
			timeout:       reportTestTimeout,
		},
		{
			name:          "namespace-memory-request",
			queryName:     "namespace-memory-request",
			newReportFunc: testFramework.NewSimpleReport,
			timeout:       reportTestTimeout + time.Minute,
		},
		{
			name:          "namespace-persistentvolumeclaim-request",
			queryName:     "namespace-persistentvolumeclaim-request",
			newReportFunc: testFramework.NewSimpleReport,
			timeout:       reportTestTimeout + time.Minute,
		},
		{
			name:          "namespace-memory-usage",
			queryName:     "namespace-memory-usage",
			newReportFunc: testFramework.NewSimpleReport,
			timeout:       reportTestTimeout + time.Minute,
		},
		{
			name:          "pod-cpu-request",
			queryName:     "pod-cpu-request",
			newReportFunc: testFramework.NewSimpleReport,
			timeout:       reportTestTimeout,
		},
		{
			name:          "pod-cpu-usage",
			queryName:     "pod-cpu-usage",
			newReportFunc: testFramework.NewSimpleReport,
			timeout:       reportTestTimeout,
		},
		{
			name:          "pod-memory-request",
			queryName:     "pod-memory-request",
			newReportFunc: testFramework.NewSimpleReport,
			timeout:       reportTestTimeout,
		},
		{
			name:          "pod-memory-usage",
			queryName:     "pod-memory-usage",
			newReportFunc: testFramework.NewSimpleReport,
			timeout:       reportTestTimeout,
		},
		{
			name:          "pod-memory-request-vs-node-memory-allocatable",
			queryName:     "pod-memory-request-vs-node-memory-allocatable",
			newReportFunc: testFramework.NewSimpleReport,
			timeout:       reportTestTimeout + time.Minute,
		},
		{
			name:          "persistentvolumeclaim-request",
			queryName:     "persistentvolumeclaim-request",
			newReportFunc: testFramework.NewSimpleReport,
			timeout:       reportTestTimeout + time.Minute,
		},
		{
			name:          "node-cpu-utilization",
			queryName:     "node-cpu-utilization",
			newReportFunc: testFramework.NewSimpleReport,
			timeout:       reportTestTimeout,
		},
		{
			name:          "node-memory-utilization",
			queryName:     "node-memory-utilization",
			newReportFunc: testFramework.NewSimpleReport,
			timeout:       reportTestTimeout,
		},
		{
			name:          "pod-cpu-request-aws",
			queryName:     "pod-cpu-request-aws",
			newReportFunc: testFramework.NewSimpleReport,
			timeout:       reportTestTimeout,
			skip:          !runAWSBillingTests,
		},
		{
			name:          "pod-memory-request-aws",
			queryName:     "pod-memory-request-aws",
			newReportFunc: testFramework.NewSimpleReport,
			timeout:       reportTestTimeout,
			skip:          !runAWSBillingTests,
		},
		{
			name:          "cluster-persistentvolumeclaim-request",
			queryName:     "cluster-persistentvolumeclaim-request",
			newReportFunc: testFramework.NewSimpleReport,
			timeout:       reportTestTimeout,
		},
		{
			name:          "cluster-cpu-capacity",
			queryName:     "cluster-cpu-capacity",
			newReportFunc: testFramework.NewSimpleReport,
			timeout:       reportTestTimeout,
		},
		{
			name:          "cluster-memory-capacity",
			queryName:     "cluster-memory-capacity",
			newReportFunc: testFramework.NewSimpleReport,
			timeout:       reportTestTimeout,
		},
		{
			name:          "cluster-cpu-usage",
			queryName:     "cluster-cpu-usage",
			newReportFunc: testFramework.NewSimpleReport,
			timeout:       reportTestTimeout,
		},
		{
			name:          "cluster-memory-usage",
			queryName:     "cluster-memory-usage",
			newReportFunc: testFramework.NewSimpleReport,
			timeout:       reportTestTimeout,
		},
		{
			name:          "cluster-cpu-utilization",
			queryName:     "cluster-cpu-utilization",
			newReportFunc: testFramework.NewSimpleReport,
			timeout:       reportTestTimeout,
		},
		{
			name:          "cluster-memory-utilization",
			queryName:     "cluster-memory-utilization",
			newReportFunc: testFramework.NewSimpleReport,
			timeout:       reportTestTimeout,
		},
		{
			name:          "aws-ec2-cluster-cost",
			queryName:     "aws-ec2-cluster-cost",
			newReportFunc: testFramework.NewSimpleReport,
			timeout:       reportTestTimeout,
			skip:          !runAWSBillingTests,
		},
	}

	scheduledReportsProduceDataTestCases := []scheduledReportProducesDataTestCase{
		{
			name:      "namespace-cpu-request-hourly",
			queryName: "namespace-cpu-request",
			timeout:   reportTestTimeout,
		},
	}

	t.Run("TestReportingProducesResults", func(t *testing.T) {
		// validate all the ReportDataSources for our tests exist before running
		// collect
		var queries []string
		waitTimeout := time.Minute

		// We wait for all ReportDataSources before anything else since even if
		// we don't use them, the collect endpoint will attempt to collect data
		// for all ReportDataSources
		_, err := testFramework.WaitForAllMeteringReportDataSourceTables(t, time.Second*5, waitTimeout)
		require.NoError(t, err, "should not error when waiting for all ReportDataSource tables to be created")

		for _, test := range reportsProduceDataTestCases {
			if test.skip {
				continue
			}
			queries = append(queries, test.queryName)
		}
		for _, test := range scheduledReportsProduceDataTestCases {
			queries = append(queries, test.queryName)
		}

		// validate all ReportGenerationQueries and ReportDataSources that are
		// used by our test cases are initialized
		testFramework.RequireReportGenerationQueriesReady(t, queries, time.Second*5, waitTimeout)

		var periodStart, periodEnd time.Time
		var collectResp operator.CollectPromsumDataResponse
		periodStart, periodEnd, collectResp = testFramework.CollectMetricsOnce(t)
		testFramework.RequireReportDataSourcesForQueryHaveData(t, queries, collectResp)

		t.Run("TestReportsProduceData", func(t *testing.T) {
			testReportsProduceData(t, testFramework, periodStart, periodEnd, reportsProduceDataTestCases)
		})
		t.Run("TestScheduledReportsProduceData", func(t *testing.T) {
			testScheduledReportsProduceData(t, testFramework, periodStart, periodEnd, scheduledReportsProduceDataTestCases)
		})
	})
}
