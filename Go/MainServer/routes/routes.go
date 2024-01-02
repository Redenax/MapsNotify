package routes

import (
	"MainServer/geocoding"
	"context"
	"crypto/tls"
	"log"
	"time"

	routespb "cloud.google.com/go/maps/routing/apiv2/routingpb"
	"google.golang.org/genproto/googleapis/type/latlng"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

const (
	fieldMask  = "*"
	apiKey     = "AIzaSyAcmeawBhpWjHpmEan-60vWy1p-L3YYjRI"
	serverAddr = "routes.googleapis.com:443"
)

func Routing(partenza string, arrivo string) *routespb.ComputeRoutesResponse {
	config := tls.Config{}
	conn, err := grpc.Dial(serverAddr,
		grpc.WithTransportCredentials(credentials.NewTLS(&config)))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()
	client := routespb.NewRoutesClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	ctx = metadata.AppendToOutgoingContext(ctx, "X-Goog-Api-Key", apiKey)
	ctx = metadata.AppendToOutgoingContext(ctx, "X-Goog-Fieldmask", fieldMask)
	defer cancel()

	latitude, longitude := geocoding.StartUpConnection(partenza)

	latitude1, longitude1 := geocoding.StartUpConnection(arrivo)

	// create the origin using a latitude and longitude
	origin := &routespb.Waypoint{
		LocationType: &routespb.Waypoint_Location{
			Location: &routespb.Location{
				LatLng: &latlng.LatLng{
					Latitude:  latitude,
					Longitude: longitude,
				},
			},
		},
	}

	// create the destination using a latitude and longitude
	destination := &routespb.Waypoint{
		LocationType: &routespb.Waypoint_Location{
			Location: &routespb.Location{
				LatLng: &latlng.LatLng{
					Latitude:  latitude1,
					Longitude: longitude1,
				},
			},
		},
	}
	req := &routespb.ComputeRoutesRequest{
		Origin:                   origin,
		Destination:              destination,
		TravelMode:               routespb.RouteTravelMode_DRIVE,
		RoutingPreference:        routespb.RoutingPreference_TRAFFIC_AWARE_OPTIMAL,
		ComputeAlternativeRoutes: false,
		Units:                    routespb.Units_METRIC,
		RouteModifiers: &routespb.RouteModifiers{
			AvoidTolls:    false,
			AvoidHighways: false,
			AvoidFerries:  true,
		},
		PolylineQuality: routespb.PolylineQuality_OVERVIEW,
	}

	// execute rpc
	resp, err := client.ComputeRoutes(ctx, req)

	if err != nil {
		// "rpc error: code = InvalidArgument desc = Request contains an invalid
		// argument" may indicate that your project lacks access to Routes
		log.Fatal(err)
	}

	return resp
}
