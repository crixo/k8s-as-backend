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
}
