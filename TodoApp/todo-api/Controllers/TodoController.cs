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

        //https://docs.microsoft.com/it-it/aspnet/core/fundamentals/http-requests?view=aspnetcore-3.1
        private readonly IHttpClientFactory _clientFactory;

        public TodoController(ILogger<TodoController> logger, IWebHostEnvironment env, IConfiguration config, IHttpClientFactory clientFactory)
        {
            _logger = logger;
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

        [HttpGet()]
        public async Task<IEnumerable<Todo>> Get()
        {

            var request = new HttpRequestMessage(
                HttpMethod.Get,
                "http://localhost:8080/apis/k8sasbackend.com/v1/namespaces/default/todos");
            request.Headers.Add("Accept", "application/json");
            request.Headers.Add("User-Agent", "TodoApp");

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

        
    }
}
