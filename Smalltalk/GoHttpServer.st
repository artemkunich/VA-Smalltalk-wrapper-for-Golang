
Object subclass: #GoHttpServer
    classInstanceVariableNames: ''
    instanceVariableNames: ''
    classVariableNames: 'Lib '
    poolDictionaries: ''!

!GoHttpServer class publicMethods !

index

	^ 'Wellcome to the Smalltalk HTTP Server'!

startOn: port
	self startOn: port logFile: nil!

startOn: port logFile: log
	
	| tmp |
	
	(tmp := self initSO) isAbtError ifTrue: [ ^ tmp ].
	
	log isNil ifFalse: [
		self startLogging: log.
	].
	
	(tmp := self runAt: port) ~= 0 ifTrue: [ ^ tmp ].
		
	[
		| rBuffer mBuffer cBuffer |
		
		rBuffer := OSStringZ new: 100.
		mBuffer := OSStringZ new: 100.
		cBuffer := OSStringZ new: 20000.
		
		[ true ] whileTrue: [
			| response resource flags |
			
			rBuffer at: 1 put: 0.
			mBuffer at: 1 put: 0.
			cBuffer at: 1 put: 0.
			
			flags := self getRequest4Resource: rBuffer method: mBuffer content: cBuffer. 
			
			flags < 0 ifFalse: [ "if flags < 0 Go HTTP Server responds automatic"
				resource := Smalltalk at: rBuffer asString asSymbol ifAbsent: [ nil ].
				(flags & 2r1) = 2r1 ifTrue: [ resource := resource new ].		
				(flags & 2r10) =  2r10 
					ifTrue: [ response := resource perform: mBuffer asString asSymbol with: cBuffer asString ] 
					ifFalse: [ response := resource perform: mBuffer asString asSymbol ].

				response isErrMessage ifTrue: [
					self putResponse: response asString asPSZ status: 500.
				] ifFalse: [
					self putResponse: response asString asPSZ.
				]
			]
		]
	] forkAt: 3! !

!GoHttpServer class privateMethods !

getRequest4Resource: rBuffer method: mBuffer content: cBuffer
	
	^(Lib at: #GetRequest) coroutineCallWith: rBuffer with: mBuffer with: cBuffer with: cBuffer size.!

initSO

	| tmp fc lb |
	
	Lib := IdentityDictionary new.
	
	lb := 'GOHTTPSERVER'.
	fc := 'RunGoHttpServer'.
	
	tmp := PlatformFunction
			callingConvention: 'c'
			function: fc
			library: lb
			parameterTypes: #(int32)
			returnType: 'int32'.
	(tmp abtProcedureAddress) isAbtError ifTrue: [ ^tmp ].	
	Lib at: fc asSymbol put: tmp.

	fc := 'PutResponse'.
	tmp := PlatformFunction
			callingConvention: 'c'
			function: fc
			library: lb
			parameterTypes: #( int32 pointer )
			returnType: 'int32'.
	(tmp abtProcedureAddress) isAbtError ifTrue: [ ^tmp ].		
	Lib at: fc asSymbol put: tmp.
	
	fc := 'GetRequest'.	
	tmp := PlatformFunction
			callingConvention: 'c'
			function: fc
			library: lb
			parameterTypes: #( pointer pointer pointer int32 )
			returnType: 'int32'.
	(tmp abtProcedureAddress) isAbtError ifTrue: [ ^tmp ].	
	Lib at: fc asSymbol put: tmp.
	
	fc := 'StartLogging'.	
	tmp := PlatformFunction
			callingConvention: 'c'
			function: fc
			library: lb
			parameterTypes: #(pointer)
			returnType: 'int32'.
	(tmp abtProcedureAddress) isAbtError ifTrue: [ ^tmp ].	
	Lib at: fc asSymbol put: tmp.
	
	fc := 'StopLogging'.	
	tmp := PlatformFunction
			callingConvention: 'c'
			function: fc
			library: lb
			parameterTypes: #()
			returnType: 'int32'.
	(tmp abtProcedureAddress) isAbtError ifTrue: [ ^tmp ].	
	Lib at: fc asSymbol put: tmp.
	!

putResponse: res

	^(Lib at: #PutResponse) coroutineCallWith: 200 with: res asPSZ!

putResponse: res status: aStatus

	^(Lib at: #PutResponse) coroutineCallWith: aStatus with: res asPSZ!

runAt: port
	
	^(Lib at: #RunGoHttpServer) coroutineCallWith: port!

startLogging: aFile

	^(Lib at: #StartLogging) coroutineCallWith: aFile asPSZ!

stopLogging

	^(Lib at: #StopLogging) coroutineCall! !

GoHttpServer initializeAfterLoad!