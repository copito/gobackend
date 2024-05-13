package gateway

// import (
// 	"go.uber.org/cadence/.gen/go/cadence/workflowserviceclient"
// 	"go.uber.org/cadence/compatibility"
// 	"go.uber.org/cadence/worker"

// 	"github.com/uber-go/tally"
// 	apiv1 "github.com/uber/cadence-idl/go/proto/api/v1"
// 	"go.uber.org/yarpc"
// 	"go.uber.org/yarpc/transport/grpc"
// 	"go.uber.org/zap"
// 	"go.uber.org/zap/zapcore"
// )

// var (
// 	HostPort       = "127.0.0.1:7933"
// 	Domain         = "test-domain"
// 	TaskListName   = "test-worker"
// 	ClientName     = "test-worker"
// 	CadenceService = "cadence-frontend"
// )

// func buildCadenceClient() workflowserviceclient.Interface {
// 	dispatcher := yarpc.NewDispatcher(
// 		yarpc.Config{
// 			Name: ClientName,
// 			Outbounds: yarpc.Outbounds{
// 				CadenceService: {
// 					Unary: grpc.NewTransport().NewSingleOutbound(HostPort),
// 				},
// 			},
// 		},
// 	)
// 	if err := dispatcher.Start(); err != nil {
// 		panic("Failed to start dispatcher")
// 	}

// 	clientConfig := dispatcher.ClientConfig(CadenceService)

// 	return compatibility.NewThrift2ProtoAdapter(
// 		apiv1.NewDomainAPIYARPCClient(clientConfig),
// 		apiv1.NewWorkflowAPIYARPCClient(clientConfig),
// 		apiv1.NewWorkerAPIYARPCClient(clientConfig),
// 		apiv1.NewVisibilityAPIYARPCClient(clientConfig),
// 	)
// }

// func startWorker(logger *zap.Logger, service workflowserviceclient.Interface) {
// 	// TaskListName identifies set of client workflows, activities, and workers.
// 	// It could be your group or client or application name.
// 	workerOptions := worker.Options{
// 		Logger:       logger,
// 		MetricsScope: tally.NewTestScope(TaskListName, map[string]string{}),
// 	}

// 	worker := worker.New(
// 		service,
// 		Domain,
// 		TaskListName,
// 		workerOptions)
// 	err := worker.Start()
// 	if err != nil {
// 		panic("Failed to start worker")
// 	}

// 	logger.Info("Started Worker.", zap.String("worker", TaskListName))
// }

// func SetupCadence() {
// 	startWorker(&zap.Logger{}, buildCadenceClient())
// 	select {}
// }
