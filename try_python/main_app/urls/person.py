from django.conf.urls import url

from main_app.views import person

urlpatterns = [
    url(r'^$', person.list, name='user-list'),
    url(r'^/(\d+)$', person.createOrModify, name='person-details'),
    url(r'^/search$', user.find, name='person-search')
]