"""admin_page URL Configuration"""
from django.conf.urls import url, include
from django.http import HttpResponseRedirect
from django.contrib import admin

from main_app.views import index

urlpatterns = [
    url(r'^admin/', admin.site.urls),
    url(r'^$', lambda r: HttpResponseRedirect('app/')),
    url(r'^app/$', index.index),
    url(r'^api/users', include('main_app.urls.person')),
]
