using System;

namespace TodoApi
{
    public class Todo
    {
        public Guid Id { get; set; }

        public string Code { get; set; }

        public DateTime When { get; set; }

        public string Message { get; set; }
    }

    public class ValidateTodoRequest
    {
        public string Raw { get; set; }
    }    

    public class ValidateTodoResponse
    {
        public bool Valid { get; set; }

        public string Message { get; set; }
    }       

    public class TodoNotification
    {
        public Guid TodoId { get; set; }

        public DateTime IssuedAt { get; set; }
    }

}
