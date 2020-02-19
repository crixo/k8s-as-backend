using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Hosting;
using Microsoft.AspNetCore.HttpsPolicy;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.Logging;
using Microsoft.OpenApi.Models;
using Serilog;

namespace TodoApi
{
    public class Startup
    {
        public Startup(IConfiguration configuration)
        {
            Configuration = configuration;
        }

        public IConfiguration Configuration { get; }

        // This method gets called by the runtime. Use this method to add services to the container.
        public void ConfigureServices(IServiceCollection services)
        {
            services.AddControllers();

            // Register the Swagger generator, defining 1 or more Swagger documents
            services.AddSwaggerGen(c =>
            {
                c.SwaggerDoc("v1", new OpenApiInfo { Title = "My API", Version = "v1" });
            });            
        }

        // This method gets called by the runtime. Use this method to configure the HTTP request pipeline.
        public void Configure(IApplicationBuilder app, IWebHostEnvironment env)
        {
            if (env.IsDevelopment())
            {
                app.UseDeveloperExceptionPage();
            }

            //app.UseHttpsRedirection();

            var routePrefix = Configuration["RoutePrefix"] ?? "";
            var fullBasePath = Configuration["FullBasePath"] ?? "";
            var relativeBasePath = Configuration["RelativeBasePath"] ?? "";
            Log.Debug("fullBasePath: {0}", fullBasePath);
            Log.Debug("relativeBasePath: {0}", relativeBasePath);

            var useSwagger = Configuration["UseSwagger"] ?? "0";

            if(useSwagger=="1"){
                // Enable middleware to serve generated Swagger as a JSON endpoint
                // https://github.com/domaindrivendev/Swashbuckle.AspNetCore/issues/1173
                app.UseSwagger(c =>
                    {
                        // c.RouteTemplate = basePath + "/swagger/{documentName}/swagger.json";
                        c.PreSerializeFilters.Add((swaggerDoc, httpReq) =>
                        {
                            swaggerDoc.Servers = new List<OpenApiServer> { 
                                new OpenApiServer { 
                                    Url = $"{httpReq.Scheme}://{httpReq.Host.Value}/{relativeBasePath}" 
                                    } };
                        });
                    }
                );

                // Enable middleware to serve swagger-ui (HTML, JS, CSS, etc.),
                // specifying the Swagger JSON endpoint.
                app.UseSwaggerUI(c =>
                {
                    c.SwaggerEndpoint(fullBasePath  + "/swagger/v1/swagger.json", "My API V1");
                    c.RoutePrefix = routePrefix;
                });            
            }

            app.UseRouting();

            app.UseAuthorization();

            app.UseEndpoints(endpoints =>
            {
                endpoints.MapControllers();
               // endpoints.MapControllerRoute("default", basePath + "/{controller=Home}/{action=Index}/{id?}");
            });
        }
    }
}
