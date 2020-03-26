using System;
using System.Collections.Generic;
using System.Linq;
using System.Net.Http;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Hosting;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.Logging;
using k8s.Models;
using System.Text;

namespace TodoApi.Controllers
{
    [ApiController]
    [Route("api/[controller]")]
    public class TodoController : ControllerBase
    {
        private static readonly string[] Summaries = new[]
        {
            "Freezing", "Bracing", "Chilly", "Cool", "Mild", "Warm", "Balmy", "Hot", "Sweltering", "Scorching"
        };

        // private static readonly HttpClient client = new HttpClient();
        //             client.DefaultRequestHeaders.Accept.Clear();
        //     // client.DefaultRequestHeaders.Accept.Add(
        //     // new MediaTypeWithQualityHeaderValue("application/vnd.github.v3+json"));
        //     client.DefaultRequestHeaders.Add("User-Agent", ".NET Foundation Repository Reporter");

        private readonly ILogger<TodoController> _logger;
        private readonly IConfiguration _config;
    //https://docs.microsoft.com/it-it/aspnet/core/fundamentals/http-requests?view=aspnetcore-3.1
        private readonly IHttpClientFactory _clientFactory;

        public TodoController(ILogger<TodoController> logger, IWebHostEnvironment env, IConfiguration config, IHttpClientFactory clientFactory)
        {
            _logger = logger;
            this._config = config;
            _logger.LogDebug(env.EnvironmentName);
            foreach(var kv in config.AsEnumerable()){
                _logger.LogDebug("k:{0} - v:{1}", kv.Key, kv.Value);
            }
            _clientFactory = clientFactory;
        }

        [HttpGet("fakes")]
        public IEnumerable<Todo> GetFakes()
        {
            var rng = new Random();
            return Enumerable.Range(1, 5).Select(index => new Todo
            {
                When = DateTime.Now.AddDays(index),
                Code = Summaries[index],
                Message = string.Format("message number {0}", rng.Next(Summaries.Length))
            })
            .ToArray();
        }

        [HttpPost("validate")]
        public async Task<ActionResult<ValidateTodoResponse>> ValidateTodo(ValidateTodoRequest dto)
        {
             _logger.LogDebug(dto.Raw);
            return new ValidateTodoResponse{Valid=true};

        }


        [HttpPost()]
        public async Task<ActionResult<Todo>> CreateTodo(Todo dto)
        {
            var ns = this._config["NAMESPACE"];

            var crd = new k8s.Models.Todo{
                ApiVersion = "k8sasbackend.com/v1",
                Kind = "Todo",
                Metadata = new V1ObjectMeta{
                    Name = dto.Code,
                    NamespaceProperty = ns,
                },
                Spec = new TodoSpec{
                    Message = dto.Message,
                    When = dto.When.ToUniversalTime()
                }
            };


            var client = _clientFactory.CreateClient();

            var svc = new k8s.TodoService(client);

            var json = svc.Serialize(crd);

            using var stringContent = new StringContent(json, Encoding.UTF8, "application/json");

            
            //TODO: make namespace as env variable
            var url = string.Format("http://localhost:8080/apis/k8sasbackend.com/v1/namespaces/{0}/todos", ns);
            _logger.LogDebug("posting todo resource to {0}", url);
            var request = new HttpRequestMessage(
                HttpMethod.Post,
                url);
            request.Headers.Add("Accept", "application/json");
            request.Headers.Add("User-Agent", "TodoApp");
            //AddBearerToken(request);
            request.Content = stringContent;

            var response = await client.SendAsync(request);
            var responseContent = await response.Content.ReadAsStringAsync();
            _logger.LogDebug(responseContent);
            if (response.IsSuccessStatusCode)
            {
                var todoCrd=svc.Convert<k8s.Models.Todo>(responseContent);

                return new Todo{
                    Id = new Guid(todoCrd.Metadata.Uid),
                    Message = todoCrd.Spec.Message, 
                    When = todoCrd.Spec.When.Value, 
                    Code=todoCrd.Metadata.Name};
            }
            else
            {
                 var status=svc.Convert<k8s.Models.V1Status>(responseContent);
                 return BadRequest(new{
                     message = status.Message,
                     code = status.Code,
                     status = status.Status,
                 });
            }
        }

        [HttpGet()]
        public async Task<IEnumerable<Todo>> GetList()
        {

            var ns = this._config["NAMESPACE"];
            var request = new HttpRequestMessage(
                HttpMethod.Get,
                string.Format("http://localhost:8080/apis/k8sasbackend.com/v1/namespaces/{0}/todos", ns));
            request.Headers.Add("Accept", "application/json");
            request.Headers.Add("User-Agent", "TodoApp");
            //AddBearerToken(request);

            var client = _clientFactory.CreateClient();

            var response = await client.SendAsync(request);

            Todo[] todos;
            if (response.IsSuccessStatusCode)
            {
                var responseContent = await response.Content.ReadAsStringAsync();
                _logger.LogDebug(responseContent);
                var svc = new k8s.TodoService(client);
                var todoList = svc.Convert<TodoList>(responseContent);
                var list = new List<Todo>();
                foreach(var item in todoList.Items){
                    list.Add(new Todo{
                        Id = new Guid(item.Metadata.Uid),
                        Code = item.Metadata.Name,
                        When = item.Spec.When.Value,
                        Message = item.Spec.Message
                    });
                }
                todos = list.ToArray();
            }
            else
            {
                todos = Array.Empty<Todo>();
            }            



            //var msg = await stringTask;
            

            return todos;
            // var rng = new Random();
            // return Enumerable.Range(1, 5).Select(index => new Todo
            // {
            //     When = DateTime.Now.AddDays(index),
            //     Code = Summaries[index],
            //     Message = string.Format("message number {0}", rng.Next(Summaries.Length))
            // })
            // .ToArray();
        }       

        private void AddBearerToken(HttpRequestMessage request)
        {
            string path = "/var/run/secrets/kubernetes.io/serviceaccount/token";
            if (System.IO.File.Exists(path))
            {
                request.Headers.Authorization = new System.Net.Http.Headers.AuthenticationHeaderValue(
                    "Bearer",
                    System.IO.File.ReadAllText(path));
            }
        }

        
    }
}
