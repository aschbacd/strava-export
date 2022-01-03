# DetailedActivity

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **int64** | The unique identifier of the activity | [optional] [default to null]
**ExternalId** | **string** | The identifier provided at upload time | [optional] [default to null]
**UploadId** | **int64** | The identifier of the upload that resulted in this activity | [optional] [default to null]
**Athlete** | [***MetaAthlete**](MetaAthlete.md) |  | [optional] [default to null]
**Name** | **string** | The name of the activity | [optional] [default to null]
**Distance** | **float32** | The activity&#39;s distance, in meters | [optional] [default to null]
**MovingTime** | **int32** | The activity&#39;s moving time, in seconds | [optional] [default to null]
**ElapsedTime** | **int32** | The activity&#39;s elapsed time, in seconds | [optional] [default to null]
**TotalElevationGain** | **float32** | The activity&#39;s total elevation gain. | [optional] [default to null]
**ElevHigh** | **float32** | The activity&#39;s highest elevation, in meters | [optional] [default to null]
**ElevLow** | **float32** | The activity&#39;s lowest elevation, in meters | [optional] [default to null]
**Type_** | [***ActivityType**](ActivityType.md) |  | [optional] [default to null]
**StartDate** | [**time.Time**](time.Time.md) | The time at which the activity was started. | [optional] [default to null]
**StartDateLocal** | [**time.Time**](time.Time.md) | The time at which the activity was started in the local timezone. | [optional] [default to null]
**Timezone** | **string** | The timezone of the activity | [optional] [default to null]
**StartLatlng** | [***LatLng**](LatLng.md) |  | [optional] [default to null]
**EndLatlng** | [***LatLng**](LatLng.md) |  | [optional] [default to null]
**AchievementCount** | **int32** | The number of achievements gained during this activity | [optional] [default to null]
**KudosCount** | **int32** | The number of kudos given for this activity | [optional] [default to null]
**CommentCount** | **int32** | The number of comments for this activity | [optional] [default to null]
**AthleteCount** | **int32** | The number of athletes for taking part in a group activity | [optional] [default to null]
**PhotoCount** | **int32** | The number of Instagram photos for this activity | [optional] [default to null]
**TotalPhotoCount** | **int32** | The number of Instagram and Strava photos for this activity | [optional] [default to null]
**Map_** | [***PolylineMap**](PolylineMap.md) |  | [optional] [default to null]
**Trainer** | **bool** | Whether this activity was recorded on a training machine | [optional] [default to null]
**Commute** | **bool** | Whether this activity is a commute | [optional] [default to null]
**Manual** | **bool** | Whether this activity was created manually | [optional] [default to null]
**Private** | **bool** | Whether this activity is private | [optional] [default to null]
**Flagged** | **bool** | Whether this activity is flagged | [optional] [default to null]
**WorkoutType** | **int32** | The activity&#39;s workout type | [optional] [default to null]
**UploadIdStr** | **string** | The unique identifier of the upload in string format | [optional] [default to null]
**AverageSpeed** | **float32** | The activity&#39;s average speed, in meters per second | [optional] [default to null]
**MaxSpeed** | **float32** | The activity&#39;s max speed, in meters per second | [optional] [default to null]
**HasKudoed** | **bool** | Whether the logged-in athlete has kudoed this activity | [optional] [default to null]
**HideFromHome** | **bool** | Whether the activity is muted | [optional] [default to null]
**GearId** | **string** | The id of the gear for the activity | [optional] [default to null]
**Kilojoules** | **float32** | The total work done in kilojoules during this activity. Rides only | [optional] [default to null]
**AverageWatts** | **float32** | Average power output in watts during this activity. Rides only | [optional] [default to null]
**DeviceWatts** | **bool** | Whether the watts are from a power meter, false if estimated | [optional] [default to null]
**MaxWatts** | **int32** | Rides with power meter data only | [optional] [default to null]
**WeightedAverageWatts** | **int32** | Similar to Normalized Power. Rides with power meter data only | [optional] [default to null]
**Description** | **string** | The description of the activity | [optional] [default to null]
**Photos** | [***PhotosSummary**](PhotosSummary.md) |  | [optional] [default to null]
**Gear** | [***SummaryGear**](SummaryGear.md) |  | [optional] [default to null]
**Calories** | **float32** | The number of kilocalories consumed during this activity | [optional] [default to null]
**SegmentEfforts** | [**[]DetailedSegmentEffort**](DetailedSegmentEffort.md) |  | [optional] [default to null]
**DeviceName** | **string** | The name of the device used to record the activity | [optional] [default to null]
**EmbedToken** | **string** | The token used to embed a Strava activity | [optional] [default to null]
**SplitsMetric** | [**[]Split**](Split.md) | The splits of this activity in metric units (for runs) | [optional] [default to null]
**SplitsStandard** | [**[]Split**](Split.md) | The splits of this activity in imperial units (for runs) | [optional] [default to null]
**Laps** | [**[]Lap**](Lap.md) |  | [optional] [default to null]
**BestEfforts** | [**[]DetailedSegmentEffort**](DetailedSegmentEffort.md) |  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

