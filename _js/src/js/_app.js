angular.module('app', [
    'ngAnimate',
    'ngResource',
    // 'ui.bootstrap',
    'ngRoute',
    'angular-loading-bar',
    'cfp.loadingBar'
]).config(['cfpLoadingBarProvider', function(cfpLoadingBarProvider) {
    console.log(cfpLoadingBarProvider);
}]);