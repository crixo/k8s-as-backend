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
using System.Collections.Concurrent;

namespace TodoApi.Controllers
{
    [ApiController]
    [Route("api/[controller]")]
    public class TodoNotificationsController : ControllerBase
    {
        private readonly ILogger<TodoController> _logger;

        private static ConcurrentDictionary<Guid, TodoNotification> TodoNotifications = 
            new ConcurrentDictionary<Guid, TodoNotification>();

        public TodoNotificationsController(ILogger<TodoController> logger)
        {
            _logger = logger;
        }   

        [HttpPost()]
        public async Task<ActionResult> AddTodoNotification(TodoNotification dto)
        {
            TodoNotifications.AddOrUpdate(dto.TodoId, dto, (key, oldValue) => dto);
            return Ok();
        }

        [HttpGet()]
        public async Task<IEnumerable<TodoNotification>> GetTodoNotifications(){
            return TodoNotifications.Values;
        }  
    }
}
