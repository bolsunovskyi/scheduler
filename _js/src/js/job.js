angular.module('app').controller('job', [
    '$scope',
    function ($scope) {

        $scope.parametrized = false;
        $scope.paramChange = function () {
            console.log($scope.parametrized);
        };


        $scope.params = [];

        $scope.addChoiceParam = function(){
            $scope.params.push({
                type: 'choice'
            });

            return false;
        };

        $scope.addStringParam = function(){
            $scope.params.push({
                type: 'string'
            });

            return false;
        };
    }
]);