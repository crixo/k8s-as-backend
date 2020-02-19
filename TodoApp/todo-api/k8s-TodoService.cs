using System.Collections.Generic;
using System.Net.Http;
using Newtonsoft.Json;
using Microsoft.Rest.Serialization;

namespace k8s
{
  public class TodoService
  {
    /// <summary>
    /// Gets or sets json serialization settings.
    /// </summary>
    public JsonSerializerSettings SerializationSettings { get; private set; }

    /// <summary>
    /// Gets or sets json deserialization settings.
    /// </summary>
    public JsonSerializerSettings DeserializationSettings { get; private set; }

    public TodoService(HttpClient httpClient) //: base(httpClient, disposeHttpClient)
    {
        Initialize();
    }

    private void Initialize()
    {
        //BaseUri = new System.Uri("http://localhost");
        SerializationSettings = new JsonSerializerSettings
        {
            Formatting = Newtonsoft.Json.Formatting.Indented,
            DateFormatHandling = Newtonsoft.Json.DateFormatHandling.IsoDateFormat,
            DateTimeZoneHandling = Newtonsoft.Json.DateTimeZoneHandling.Utc,
            NullValueHandling = Newtonsoft.Json.NullValueHandling.Ignore,
            ReferenceLoopHandling = Newtonsoft.Json.ReferenceLoopHandling.Serialize,
            ContractResolver = new ReadOnlyJsonContractResolver(),
            Converters = new  List<JsonConverter>
                {
                    new Iso8601TimeSpanConverter()
                }
        };
        DeserializationSettings = new JsonSerializerSettings
        {
            DateFormatHandling = Newtonsoft.Json.DateFormatHandling.IsoDateFormat,
            DateTimeZoneHandling = Newtonsoft.Json.DateTimeZoneHandling.Utc,
            NullValueHandling = Newtonsoft.Json.NullValueHandling.Ignore,
            ReferenceLoopHandling = Newtonsoft.Json.ReferenceLoopHandling.Serialize,
            ContractResolver = new ReadOnlyJsonContractResolver(),
            Converters = new List<JsonConverter>
                {
                    new Iso8601TimeSpanConverter()
                }
        };
    }    

    public T Convert<T>(string rawJsonContent)
    {
      return SafeJsonConvert.DeserializeObject<T>(rawJsonContent, DeserializationSettings);
    }     

    public string Serialize(object obj)
    {
      return SafeJsonConvert.SerializeObject(obj, SerializationSettings);
    } 
  }

}



