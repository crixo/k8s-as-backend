using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Hosting;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.Logging;

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

        private readonly ILogger<TodoController> _logger;

        public TodoController(ILogger<TodoController> logger, IWebHostEnvironment env, IConfiguration config)
        {
            _logger = logger;
            _logger.LogDebug(env.EnvironmentName);
            foreach(var kv in config.AsEnumerable()){
                _logger.LogDebug("k:{0} - v:{1}", kv.Key, kv.Value);
            }

        }

        [HttpGet]
        public IEnumerable<Todo> Get()
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
    }
}
