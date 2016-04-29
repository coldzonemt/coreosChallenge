angular.module('app', [
  'ngRoute', //use ui-route instead of ngRoute
  //include other modules to be used here
])
//config is run before the application is spun up. Routes need to be in place first. 
.config(function($routeProvider, $locationProvider) {
  'use strict';
  $locationProvider.html5Mode(true);
  $routeProvider
    .when('/', {
      controller: 'LoginPageCtrl',
      templateUrl: '/login.html',
      title: 'Login',
    })
    .when('/list', {
      controller: 'ListPageCtrl',
      templateUrl: '/list.html',
      title: 'List',
    });
    //consider using "$urlRouterProvider.otherwise('/')" here for any other address that is typed in
})
.controller('LoginPageCtrl', function($scope, $http, $location, $rootScope) {
  $scope.user = {
    username: '',
    password: '',
  };

  $scope.login = function() {
    $http.post('/api/login', $scope.user)
    .then(function(resp) {
      $rootScope.token = resp.data.token;
      $location.path('/list');
    }, function(e) {
      alert ('error ', e);
    });
  };
})
.controller('ListPageCtrl', function($scope, $http, $rootScope, $location) {
  if (!$rootScope.token) {
    $location.path('/');
    return;
  }

  $scope.ghUsername = '';
  $scope.items = [];
  $scope.total = 0;

  $scope.submit = function() {
    $http.put('/api/list?token=' + $rootScope.token + '&gh-username=' + $scope.ghUsername)
    .then(function(resp) {
      for (i=0; i < resp.data.length; i++) {
        $scope.total += resp.data[i].stargazers_count;
      }
      $scope.items = resp.data;
      $scope.items.sort(function (a, b) {
        if (a.stargazers_count > b.stargazers_count) {
          return -1;
        }
        if (a.stargazers_count < b.stargazers_count) {
          return 1;
        }
        // a must be equal to b
        return 0;
      });
    });
  };

});
