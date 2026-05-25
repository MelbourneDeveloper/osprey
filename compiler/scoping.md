Structured Handling of Scoped Eﬀects
Zhixuan Yang1, Marco Paviotti1, Nicolas Wu1, Birthe van den Berg2, and Tom
Schrijvers2
1 Imperial College London, United Kingdom
2 KU Leuven, Belgium
Abstract. Algebraic eﬀects oﬀer a versatile framework that covers a
wide variety of eﬀects. However, an important family of operations—
those that delimit a scope, and usually modelled as handlers—are not
algebraic, thus preventing them from being used freely in conjunction
with algebraic operations. Although proposals for scoped operations ex-
ist, they are either ad-hoc and unprincipled, or too inconvenient for
practical programming. This paper provides the best of both worlds:
a theoretically-founded model of scoped eﬀects that is convenient for im-
plementation and reasoning. Our new model is based on an adjunction
between a locally finitely presentable category and a category of functo-
rial algebras. Using comparison functors between adjunctions, we show
that our new model, the earlier indexed model, and a third approach
that simulates scoped operations in terms of algebraic ones have equal
expressivity for handling scoped operations. We consider our new model
to be the sweet spot between ease of implementation and structuredness.
Additionally, our approach automatically induces fusion laws of handlers
of scoped eﬀects, which are useful for reasoning and optimisation.
Keywords: Computational eﬀects· Category theory· Haskell· Alge-
braic theories· Scoped eﬀects· Handlers· Abstract syntax
1 Introduction
For a long time monads [48, 63, 71] have been the go-to approach for purely
functional modelling of and programming with side eﬀects. However, in recent
years an alternative approach, algebraic eﬀects [51], is gaining more traction. A
big breakthrough has been the introduction of handlers [55], which has made
algebraic eﬀects suitable for programming and has led to numerous dedicated
languages and libraries implementing algebraic eﬀects and handlers. In compar-
ison to monads, algebraic eﬀects provide a more modular approach to computa-
tions with eﬀects, in which the syntax and semantics of eﬀects are separated—
computations invoking algebraic operations can be defined syntactically, and the
semantics of operations are given by handlers separately in possibly many ways.
A disadvantage of algebraic eﬀects is that they are less expressive than mon-
ads; not all eﬀects can be easily expressed or composed within their confines.
For instance, operations like catch for exception handling, spawn for parallel
2
composition of processes, or once for restricting nondeterminism are not con-
ventional algebraic operations; instead they delimit a computation within their
scope. Such operations are usually modelled as handlers, but the problem is that
they cannot be freely used amongst other algebraic operations: when a handler
implementing a scoped operation is applied to a computation, the computation
is transformed from a syntactic tree of algebraic operations into some semantic
domain implementing the scoped operation. Consequently, all subsequent oper-
ations on the computation can only be given in the particular semantic domain
rather than as mere syntactic operations, thus nullifying the crucial advantage
of modularity when separating syntax and semantics of eﬀects.
To remedy the situation, Wu et al. [73] proposed a practical, but ad-hoc,
generalization of algebraic eﬀects in Haskell that encompasses scoped eﬀects,
that has been adopted by several algebraic eﬀects libraries [59, 45, 35]. More
recently, Pir´og et al. [49] sought to put this ad-hoc approach for scoped eﬀects
on the same formal footing as algebraic eﬀects. Their solution resulted in a
construction based on a level-indexed category, called indexed algebras, as the
way to give semantics to scoped eﬀects. However, this formalisation introduces
a disparity between syntax and semantics that makes indexed algebras not as
structured as the programs they interpret, where they use an ad-hoc hybrid
fold that requires indexing for the handlers, but not for the program syntax.
Moreover, indexed algebras are not ideal for widespread implementation as they
require dependent typing, in at least a limited form like gadts [28].
This paper presents a more structured way of handling scoped eﬀects, which
we call functorial algebras. They are principled and formally grounded on cat-
egory theory, and at the same time more structured than using the indexed
algebras of Pir´og et al. [49], in the sense that the structure of functorial algebras
directly follow the abstract syntax of programs with scoped eﬀects. Additionally,
our approach can be practically implemented without the need for dependent
types or gadts, making it available for a wider range of programming languages.
In particular, we make the following contributions:
– We highlight the loss of modularity when modelling scoped operations as
handlers and sketch how the problem is solved using functorial algebras in
Haskell, along with a number of programming examples (Section 2);
– We develop a category-theoretic foundation of functorial algebras as a notion
of handlers of scoped eﬀects. Specifically, we show that there is an adjunc-
tion between functorial algebras and a base category, inducing the monad
modelling the syntax of scoped eﬀects (Section 3);
– We show that the expressivity of functorial algebras, Pir´og et al. [49]’s in-
dexed algebras, and simulating scoped eﬀects with algebraic operations and
recursion are equal, by constructing interpretation-preserving functors be-
tween the three categories of algebras (Section 4);
– We present the fusion law of functorial algebras, which directly follows from
the naturality of the adjunction underlying functorial algebras (Section 5).
Finally, we discuss related work (Section 6) and conclude (Section 7).
3
2 Scoped Eﬀects for the Working Programmer
We start with a brief recap of handlers of algebraic eﬀects (Section 2.1), and then
we highlight the loss of modularity when modelling some non-algebraic eﬀectful
operations as handlers (Section 2.2). We then show how the problem is solved by
modelling them as scoped operations and handling them with functorial algebras
in Haskell (Section 2.3), whose categorical foundation will be developed more
carefully in later sections.
2.1 Handlers of Algebraic Eﬀects
For the purpose of demonstration, in this section we base our discussion on a
simplistic implementation of eﬀect handlers in Haskell using free monads, al-
though the problem with eﬀect handlers highlighted in this section applies to
other more practical implementations of eﬀect handlers, either as libraries (e.g.
[36, 30]) or standalone languages (e.g. [7, 39, 43]).
Following Plotkin and Pretnar [55], computational eﬀects, such as exceptions,
mutable state, and nondeterminism, are described by signatures of primitive
eﬀectful operations, and they can be realised as Haskell functors. For example,
the following functor (with the evident Functor instance) is the signature of three
operations, throwing an exception, reading and writing a mutable Int-state:
data ES x= Throw |Put Int x |Get (Int →x) (1)
Typically, a constructor of a signature functor Σ has a type isomorphic to P →
(R →x) →Σ x for some types P and R. As in (1), the types of the three
constructors are isomorphic to Throw :: () →(Void →x) →ES x, Put :: Int →
(() →x) →ES x and Get :: () →(Int →x) →ES x respectively where Void is
the empty type. Each constructor of a signature functor Σ is thought of as an
operation that takes a parameter of type P and produces a result of type R, or
equivalently, has R-many possible ways to continue the computation after the
operation. Given any (signature) functor Σ, computations invoking operations
from Σ are modelled by the following datatype, called the free monad of Σ,
data Free Σ a = Return a |Call (Σ (Free Σ a))
whose first case represents a computation that just returns a value, and the
second case represents a computation calling an operation from Σ with more
Free Σ a subterms as arguments, which are understood as the continuation of
the computation after this call, depending on the outcome of this operation.
The inductive datatype Free Σ a comes with a recursion principle:
handle :: (Σ b →b) →(a →b) →Free Σ a →b
handle alg g (Return x) = g x
handle alg g (Call op) = alg (fmap (handle alg g) op)
which folds a tree of operations Free Σ a into a type b, providing a way Σ b →b,
usually called a Σ-algebra, to perform operations from Σ on b and a way a →b
4
to transform the returned type a of computations to b. The function handle can
be used to give Free Σ a monad instance:
return :: a →Free Σ a
return= Return
(>> =) :: Free Σ a →(a →Free Σ b) →Free Σ b
m >> = k= handle Op k m
The monadic instance allows the programmer to build eﬀectful computations
using the do-notation in a clean way. For example, the following program updates
the state s to n /s for some n :: Int and throwing an exception when s is 0:
safeDiv :: Int →Free ES Int
safeDiv n= do {s ←get; if s ≡0 then throw else put (n /s); return (n /s)}
where the auxiliary ‘smart constructors’ dealing with Call appropriately are
get= Call (Get Return); put n= Call (Put n (Return ())); throw= Call Throw
Clearly, the free monad merely models eﬀectful computations syntactically
without specifying how these operations are actually implemented. Indeed, the
program safeDiv above is defined without saying how mutable state and excep-
tions are implemented at all. To actually give useful semantics to programs built
with free monads, the programmer uses the handle function above to interpret
programs with Σ-algebras, which are called handlers in this context.
For example, given a program r ::Free ES a for some a, a handler catchHdl r ::
ES (Free ES) →Free ES that gives the usual semantics to throw is
catchHdl :: Free ES a →ES (Free ES a) →Free ES a
catchHdl r Throw= r; catchHdl r op= Call op (2)
which evaluates r for recovery in case of throwing an exception, and leaves
other operations untouched in the free monad. An important advantage of the
approach of eﬀect handlers is that diﬀerent semantics of a computational eﬀect
can be given by diﬀerent handlers. For example, suppose that in some scenario
one would like to interpret exceptions as unrecoverable errors and stop the exe-
cution of the program when an exception is raised. Then the following handler
can be defined for this behaviour:
catchHdl′:: Free ES a →ES (Free ES (Maybe a)) →Free ES (Maybe a)
catchHdl′r Throw= return Nothing; catchHdl′r op= Call op (3)
As expected, applying these two handlers to the program safeDiv 5 produces
diﬀerent results:
handle (catchHdl (return 42)) return (safeDiv 5) :: Free ES Int
= do {s ←get; if s ≡0 then return 42 else put (n /s); return (n /s)}
handle (catchHdl′(return 42)) (return·Just) (safeDiv 5) :: Free ES (Maybe Int)
= do {s ←get; if s ≡0 then return Nothing else put (n /s); return (Just (n /s))}
Note that exception throwing and catching are modelled diﬀerently in the ap-
proach of algebraic eﬀects and handlers, one as an operation in the signature
5
ES and one as a handler, although it is natural to expect both of them to
be operations of the eﬀect of exceptions. This asymmetry results from the fact
that exception catching is not algebraic: if catch was modelled as a binary op-
eration in the signature, then the monadic bind >> = of the free monad earlier,
which intuitively means sequential composition of programs, would imply that
(catch r p) >> = k= catch (r >> = k) (p >> = k), which is semantically undesirable.
Thus the perspective of Plotkin and Pretnar [55] is that non-algebraic operations
like catch should be deemed diﬀerent from algebraic operations, and they can
be modelled as handlers (of algebraic operations).
2.2 Scoped Operations as Handlers Are Not Modular
However, this treatment of non-algebraic operations leads to a somewhat subtle
complication: as observed by Wu et al. [73], when non-algebraic operations (such
as catch) are modelled with handlers, these handlers play a dual role of (i) mod-
elling the syntax of the operation (the scope for which exceptions are caught
by catch) and (ii) giving semantics to it (when an exception is caught, run the
recovery program). To see the problem more concretely, ideally one would like to
have a syntactic operation catch of the following type that acts on computations
without giving specific semantics a priori,
catch :: Free ES a →Free ES a →Free ES a
allowing to write programs like
prog= do {x ←catch (safeDiv 5) (return 42); put (x + 1)} (4)
and the semantics of (both algebraic and non-algebraic) operations in prog can be
given separately by handlers. Unfortunately, when catch is modelled as handlers
catchHdl or catchHdl′as in the last subsection, the program prog must be written
diﬀerently depending on which handler is used:
do x ←handle (catchHdl (return 42)) return (safeDiv 5); put (x + 1)
vs. do xMb ←handle (catchHdl′(return 42)) (return·Just) (safeDiv 5)
case xMb of {Nothing →return Nothing
(Just x) →do r ←put (x + 1); return (Just r)}
The issue is that these handlers interpret the operation catch in diﬀerent seman-
tic domains, Free ES a and Free ES (Maybe a), and this aﬀects both the value
x that is returned, and the way the subsequent put is expressed. Therefore, non-
algebraic operation catch modelled as handlers is not as modular as algebraic
operations, weakening the advantage of programming with algebraic eﬀects.
2.3 Scoped Eﬀects and Functorial Algebras
Now we present an overview of a solution to the problem highlighted above
by modelling exception catching as scoped eﬀects [49] and handle them using
functorial algebras, which will be more formally developed in later sections.
6
Syntax of Scoped Operations To achieve modular interpretation of (non-algebraic)
operations delimiting scopes, such as catch, which are called scoped operations
by Pir´og et al. [49], we need to generalise the free monad Free Σ to a monad
Prog Σ Γ accommodating both algebraic and scoped operations. The monad
is parameterised by two functors Σ and Γ, which we call the algebraic signa-
ture and the scoped signature respectively. The intention is that a constructor
(R →x) →Σ x of the algebraic signature represents an algebraic operation that
produces a result of type R as usual, while a constructor (N →x) →Γ x of
the scoped signature represents a scoped operation that creates N -many scopes
enclosing programs.
Example 1. As in the previous subsection, the eﬀect of exceptions has an alge-
braic operation for throwing exceptions, which produces no values, and a scoped
operation for catching exceptions, which creates two scopes, one enclosing the
program for which exceptions are caught, and the other enclosing the recovery
computation. Thus the algebraic and scoped signatures are respectively
data Throw x= Throw data Catch x= Catch x x (5)
Example 2. An eﬀect of explicit nondeterminism has two algebraic operations
for nondeterministic choice and a scoped operation Once:
data Choice x= Fail |Or x x data Once x= Once x (6)
The intention is that this eﬀect implements logic programming [22]—solutions
to a problem are exhaustively searched: operation Or p q splits a search branch
into two; Fail marks a failed branch; and the scoped operation Once p keeps
only the first solution found by p, making it semi-deterministic, which is useful
for speeding up the search with heuristics from the programmer.
Similar to the free monad, the Prog monad models the syntax of computa-
tions invoking operations from Σ and Γ:
data Prog Σ Γ a = Return a |Call (Σ (Prog Σ Γ a))
| Enter (Γ (Prog Σ Γ (Prog Σ Γ a))) (7)
Thus an element of Prog Σ Γ a can either (i) return an a-value without
causing eﬀects, or (ii) call an algebraic operation in Σ with more subterms
of Prog Σ Γ a as the continuation after the operation, or (iii) enter the scope of
a scoped operation. The third case deserves more explanation: the first Prog in
(Γ (Prog Σ Γ (Prog Σ Γ a))) represents the programs enclosed by the scoped
operation, and the second Prog represents the continuation of the program after
the scoped operation, and thus the boundary between programs inside and out-
side the scope is kept in the syntax tree, which is necessary because collapsing
the boundary might change the meaning of a program. The distinction between
algebraic and scoped operations can be seen more clearly from the following
monadic bind of Prog (the monadic return of Prog is still Return):
(>> =) :: Prog Σ Γ a →(a →Prog Σ Γ) b →Prog Σ Γ b
(Return a) >> = k= k a; (Call op) >> = k= Call (fmap (>> =k) op)
(Enter sc) >> = k= Enter (fmap (fmap (>> =k)) sc)
7
For algebraic operations, extending the continuation (>> =k) directly acts on the
argument to the algebraic operation, whereas for scoped operation, (>> =k) acts
on the second layer of Prog. Thus for an algebraic operation o, (o p) >> = k and
o (p >> = k) have the same representation, whereas for a scoped operation s,
(s p) >> = k and s (p >> = k) have diﬀerent representations, which is precisely the
distinction between algebraic and scoped operations.
The constructors Call and Enter are clumsy to work with, and for writing
programs more naturally, we define smart constructors for operations. Generally,
for algebraic operations Op ::F x →Σ x and scoped operations Sc ::G x →Γ x,
the smart constructors are
op :: F (Prog Σ Γ a) →Prog Σ Γ a op= Call·Op sc :: G (Prog Σ Γ a) →Prog Σ Γ a
sc= Enter·fmap (fmap return)·Sc
For example, the smart constructor for Catch (Example 1) is
catch :: Prog Σ Catch a →Prog Σ Catch a →Prog Σ Catch a
catch h r= Enter (Catch (fmap return h) (fmap return r))
With all machinery in place, now we can define the program (4) using Prog that
we could not write with Free:
prog= do {x ←catch (safeDiv 5) (return 42); put (x + 1)}
Handlers of Scoped Operations Similar to Free, the Prog monad merely models
the syntax of eﬀectful computations, and more useful semantics need to be given
by handlers. Although Pir´og et al. [49] developed a notion of indexed algebras
for this purpose, they turn out to be more complicated than necessary (we will
discuss them in Section 4), and the contribution of this paper is precisely a new
kind of handlers for scoped operations, which we call functorial algebras.
Given signatures Σ and Γ, a functorial algebra for them is a quadruple
⟨f, b, ealg, balg⟩for some functor f called the endofunctor carrier, type b called
the base carrier. The other two components ealg :: EndoAlg Σ Γ f and balg ::
BaseAlg Σ Γ f b are called the endofunctor algebra and the base algebra. Their
types are fully shown in Figure 1. The intuition is that functor f and ealg inter-
pret the part of a program enclosed by scoped operations, and the type b and
balg interpret the part of a program not enclosed by any scopes.
Example 3. The standard semantics of exception catching (cf. handler (2)) can
be implemented by a functorial algebra with the conventional Maybe functor as
the endofunctor carrier with the following EndoAlg:
excE :: EndoAlg Throw Catch Maybe
excE= EndoAlg {..}where enterE :: Catch (Maybe (Maybe a))
returnE= Just →Maybe a
callE Throw= Nothing enterE (Catch Nothing r) = join r
enterE (Catch (Just k) ) = k
For the base carrier that interprets operations not enclosed by any catch, a
straightforward choice is just taking Maybe a as the base carrier for a type a,
8
data EndoAlg Σ Γ f= EndoAlg {
returnE :: ∀x. x →f x,
callE :: ∀x. Σ (f x) →f x,
enterE :: ∀x. Γ (f (f x)) →f x }
data BaseAlg Σ Γ f a=
BaseAlg {callB :: Σ a →a
, enterB :: Γ (f a) →a }
hcata :: (Functor Σ, Functor Γ ) ⇒(EndoAlg Σ Γ f ) →Prog Σ Γ a →f a
hcata alg (Return x) = returnE alg x
hcata alg (Call op) = (callE alg·fmap (hcata alg)) op
hcata alg (Enter scope) = (enterE alg·fmap (hcata alg·fmap (hcata alg))) scope
handle :: (Functor Σ, Functor Γ )
⇒(EndoAlg Σ Γ x) →(BaseAlg Σ Γ x b) →(a →b) →Prog Σ Γ a →b
handle ealg balg gen (Return x) = gen x
handle ealg balg gen (Call op) = (callB balg·fmap (handle ealg balg gen)) op
handle ealg balg gen (Enter sc)
= (enterB balg·fmap (hcata ealg·fmap (handle ealg balg gen))) sc
Fig. 1: A Haskell implementation of handling with functorial algebras
and setting callB= callE and enterB= enterE, which means that operations
inside and outside scopes are interpreted in the same way.
In general, we can define a specialised version of handle (Figure 1) that only
takes an endofunctor algebra as input for interpreting operations inside and
outside scopes in the same way:
handleE :: (EndoAlg Σ Γ f ) →Prog Σ Γ a →f a
handleE ealg@(EndoAlg {..}) = handle ealg (BaseAlg callE enterE) returnE
Applying handleE excE to the following program produces Just 43 as expected.
do {x ←catch throw (return 42); return (x + 1)} (8)
For the non-standard semantics (cf. (3)) that disables exception recovery, one
can define another endofunctor algebra excE′by replacing enterE in exeE with
enterE′:: Catch (Maybe (Maybe a)) →Maybe a
enterE′(Catch Nothing ) = Nothing; enterE′(Catch (Just k) ) = k
With excE′, handling the program in (8) produces Nothing as expected.
Now we provide some intuition for how functorial algebras work. First note
that the three fields of EndoAlg in Figure 1 precisely correspond to the three
cases of Prog (7). Thus by replacing the constructors of Prog with the correspond
fields of EndoAlg, we have a polymorphic function hcata ealg ::∀x. Prog Σ Γ x →
f x (Figure 1) turning a program into a value in f.
The function handle (Figure 1) takes a functorial algebra, a function gen ::
a →b and a program p as arguments, and it handles all the eﬀectful operations in
p by using hcata ealg for interpreting the part of p inside scoped operations and
9
balg for interpreting the outermost layer of p outside any scoped operations. The
function gen corresponds to the ‘value case’ of handlers of algebraic eﬀects, which
transforms the a-value returned by a program into the type b for interpretation.
We close this section with some more examples of handling scoped eﬀects
with functorial algebras. The supplementary material of this paper also contains
an OCaml implementation of functorial algebras and the following examples.
Example 4. The standard way to handle explicit nondeterminism with the semi-
deterministic operator once (Example 2) is using a functorial algebra with the
list functor as the endofunctor carrier together with the following algebra:
ndetE :: EndoAlg Choice Once [ ] enterE :: Once [[a ]] →[a ]
ndetE= EndoAlg {..}where enterE (Once x) =
callE :: Choice [a ] →[a ] if x ≡[ ] then [ ] else head x
callE Fail = [ ] returnE :: a →[a ]
callE (Or x y) = x ++y returnE x = [x ]
Then applying handleE ndetE to the following program produces [1, 2] as ex-
pected. In comparison, if once were algebraic, the result would be [1].
do {n ←once (or (return 1) (return 3)); or (return n) (return (n + 1))}
Example 5. In the last example we used the list functor to interpret explicit
nondeterminism, resulting in the depth-first search (DFS) strategy for searching.
Noted by Spivey [62], other search strategies can be implemented by other choices
of functors. For example, depth-bounded search (DBS) can be implemented with
the functor Int →[a ], and breadth-first search (BFS) can be implemented with
the functor [[a ]] (or Kidney and Wu [34]’s more eﬃcient LevelT functor).
A powerful application of scoped eﬀects is that we can model search strategies
as scoped operations:
data Strategy x= DFS x |BFS x |DBS Int x
so that the programmer can freely specify the search strategy of nondeterministic
choices in a scope. The algebraic signature Choice and scoped signature Strategy
can be handled by a functorial algebra carried by endofunctor ([a ], [[a ]], Int →
[a ]) and base type [a ] (assuming depth-first search is the default strategy). The
complete code can be found in the supplementary material.
Example 6. A scoped operation for the eﬀect of mutable state is the operation
local s p that executes the program p with a state s and restores to the origi-
nal state after p finishes. Clearly, (local s p; k) is diﬀerent from local s (p; k), and
thus local should be modelled as a scoped operations of signature data Local s a=
Local s a. Together with the usual algebraic operations get and put of state,
Local can be interpreted with a functorial algebra carried by the state monad
type State s a= s →(s, a). The essential part of the functorial algebra is the
following enterE for Local (complete code in the supplementary material):
enterE :: Local (State s (State s a)) →State s a
enterE (Local s′f ) s = let ( , k) = f s in k s
10
Example 7. Parallel composition of processes is not an operation in the usual
algebraic presentations of process calculi [64, 65] precisely because it not alge-
braic: (p |q) >> = k ̸= (p >> = k) |(q >> = k). Again, we can model it as a scoped
operation, and diﬀerent scheduling behaviours of processes can be given as dif-
ferent functorial algebras. The supplementary material contains complete code
of handling parallel composition using the so-called resumption monad [11, 50]:
data Resumption m a= More (m (Resumption m a)) |Done a
3 Categorical Foundations for Scoped Operations
We now move on to a categorical foundation for scoped eﬀects and functorial
algebras. First, we recall some standard category theory underlying algebraic
eﬀects and handlers (Section 3.1) and also Pir´og et al. [49]’s monad P that
models the syntax of scoped operations, which is exactly the Prog monad in
the Haskell implementation (Section 3.2). Then, we define functorial algebras
formally (Section 3.3) and show that there is an adjunction between the category
of functorial algebras and the base category (Section 3.4) inducing the monad
P, which provides a means to interpret the syntax of scoped operations.
The rest of this paper assumes familiarity with basic category theory, such
as adjunctions, monads, and initial algebras, which are covered by standard
texts [44, 6, 58]. The mathematical notation in this paper is summarised in
Appendix B, which may be consulted if the meaning of some symbols are unclear.
3.1 Syntax and Semantics of Algebraic Operations
The relationships between equational theories, Lawvere theories, monads, and
computational eﬀects are well-studied for decades from many perspectives [60,
32, 57, 51, 26, 48]. Here we recap a simplified version of equational theories by
Kelly and Power [32] that we follow to model algebraic and scoped eﬀects on
locally finitely presentable (lfp) categories [1]. However, we only consider unen-
riched categories in this paper.
Locally Finitely Presentable Categories The use of lfp categories in this paper
is limited to some standard results about the existence of many initial algebras
in lfp categories, and thus a reader not familiar with lfp categories may follow
this paper with some simple intuition: a category C is lfp if it has all (small)
colimits and a set of finitely presentable objects such that every object in C can be
obtained by ‘glueing’ (formally, as filtered colimits of) some finitely presentable
objects. For example, Set is lfp with finite sets as its finitely presentable objects,
and indeed every set can be obtained by glueing, here meaning taking the union
of, all its finite subsets: X= {N ⊆X |N finite }. Other examples of lfp
categories include the category of partially ordered sets, the category of graphs,
the category of small categories, and presheaf categories (we refer the reader
to the excellent exposition [60] for concrete examples), thus lfp categories are
widespread to cover many semantic settings of programming languages.
11
Moreover, an endofunctor F : C →C is said to be finitary if it preserves
‘glueing’ (filtered colimits), which implies that its values FX are determined
by its values at finitely presentable objects: FX∼
= F(colimiNi)∼
= colimiFNi
where Ni are the finitely presentable objects that generate X when glued to-
gether. For example, polynomial functors n∈N Pn×(−)n on Set are finitary
where Pn is a set for every n.
Algebraic Operations on LFP Categories Fixing an lfp category C, we take fini-
tary endofunctors Σ : C →C as signatures of operations on C. Like in Sec-
tion 2.1, the intuition is that every natural transformation C(R,−) P →Σ−
for some object P : C and finitely presentable object R : C standards for an
operation taking a parameter of type P and R-many arguments. The category
Σ-Alg of Σ-algebras is defined as usual: it has pairs ⟨X : C, α : ΣX →X⟩
as objects and morphisms h : X →X′ such that h·α = α′
·Σh as morphisms
⟨X, α⟩→⟨X′
, α′⟩. The following classical results (see e.g. [2, 5]) give suﬃcient
conditions for constructing initial and free Σ-algebras:
Lemma 1. If category C has finite coproducts and colimits of all ω-chains and
functor Σ : C →C preserves them, then the forgetful functor UΣ : Σ-Alg →C
forgetting the structure maps has a left adjoint FreeΣ : C →Σ-Alg mapping
every X : C to a Σ-algebra ⟨Σ∗X, opX⟩where Σ∗X denotes the initial algebra
µY. X+ ΣY and opX : ΣΣ∗X →Σ∗X.
Lemma 1 is applicable to our setting since C being lfp directly implies that
it has all colimits, and finitary functors Σ preserve colimits of ω-chains because
colimits of ω-chains are filtered. Hence we have an adjunction: FreeΣ ⊣UΣ :
Σ-Alg →C. We denote the monad from the adjunction by Σ∗
= UΣFreeΣ
(which is implemented as the Free Σ monad in Section 2.1). The idea is still
that syntactic terms built from operations in Σ are modelled by the monad
Σ∗, and semantics of operations are given by Σ-algebras. Given any Σ-algebra
⟨X, α: ΣX →X⟩and morphism g: A→X in C, they induce an interpretation
morphism handle⟨X,α⟩g: Σ∗A→X s.t.
handle⟨X,α⟩g= UΣ(ϵ⟨X,α⟩·FreeΣg) : Σ∗A= UΣFreeΣA→X (9)
where ϵ⟨X,α⟩: FreeΣUΣ⟨X, α⟩→⟨X, α⟩is the counit of FreeΣ ⊣UΣ.
Algebraic Eﬀects and Handlers The perspective of Plotkin and Pretnar [55] is
that computational eﬀects are characterised by signatures Σof primitive eﬀectful
operations, and they determine monads Σ∗ that model programs syntactically.
Then Σ-algebras are handlers [55] of operations that can be applied to programs
using (9) to give specific semantics to operations.
The approach of algebraic eﬀects has led to a significant body of research on
programming with eﬀects and handlers, but it imposes an assumption on the
operations to be modelled: the construction of Σ∗ in Lemma 1 [2, 5] implies
that the multiplication µ of the monad Σ∗ satisfies the algebraicity property:
op·(Σ◦µ) = µ·(op ◦Σ∗) : ΣΣ∗Σ∗ → Σ∗ where op : Σ(Σ∗) → Σ∗. This
12
intuitively means that every operation in Σmust be commutative with sequential
composition of computations. Many, but not all, eﬀectful operations satisfy this
property, and they are called algebraic operations.
Adjoint Approach to Eﬀects The crux of algebraic eﬀects and handlers is the
adjunction FreeΣ ⊣UΣ. However, we have not relied on the adjunction being
the free/forgetful one at all: given any monad P : C →C that models the syntax
of eﬀectful Programs, if L⊣R : D →C is an adjunction such that RL∼
= P as
monads, then objects D in D provide a means to interpret programs PA—for
any g: A→RD in C, we have the following interpretation morphism
handleDg= R(ϵD·Lg) : PA∼
= R(LA) →RD (10)
The intuition for g is that it transforms the returned value A of a computation
into the carrier RD, so it corresponds to the ‘value case’ of eﬀect handlers [8].
Pir´og et al. [49] call this approach the adjoint-theoretic approach to syntax and
semantics of eﬀects, and they construct an adjunction between indexed algebras
and the base category for modelling scoped operations. Earlier, Levy [40] and
Kammar and Plotkin [31] also adopt a similar adjunction-based viewpoint in the
treatment of call-by-push-value calculi: value types are interpreted in the base
category C, and computation types are interpreted in the algebra category D.
Remark 1. A notable missing part of our treatment is the equations that specify
operations in a signature. Following Kelly and Power [32], an equation for a signa-
ture Σ : C →C can be formulated as a pair of monad morphisms σ, τ : Γ∗→Σ∗
τ
for some finitary functor Γ, and taking their coequaliser Γ∗ Σ∗ M
in
σ
the category of finitary monads constructs a monad M that represents terms
modulo the equation l= r. Although it seems straightforward to extend this
formulation of equational theories work with scoped eﬀects, we do not consider
equations in this paper for the sake of simplicity.
Remark 2. Working with lfp categories precludes operations with infinite argu-
ments, such as the get operation (1) of mutable state when the state has infinite
possible values, but this limitation is not inherent and can be handled by moving
to locally κ-presentable categories [1] for some larger cardinal κ.
3.2 Syntax of Scoped Operations
Not all operations in programming languages can be adequately modelled as
algebraic operations on Set, for example, λ-abstraction [17], memory cell gener-
ation [51, 41], more generally, eﬀects with dynamically generated instances [65],
explicit substitution [20], channel restriction in π-calculus [64], and their syntax
are usually modelled in some functor categories. More recently, Pir´og et al. [49]
extend Ghani et al. [20]’s work to model a family of non-algebraic operations,
which they call scoped operations. In this subsection, we review their develop-
ment in the setting of lfp categories. Throughout the rest of the paper, we fix an
13
Γ(PX) (11)
lfp category C, and refer to it as the base category, and it is intended to be the
category in which types of a programming language are interpreted. Further-
more, we fix two finitary endofunctors Σ, Γ : C →C and call them the algebraic
signature and scoped signature respectively.
Syntax Endofunctor P Now our goal is to construct a monad P : C →C that
models the syntax of programs with algebraic operations in Σ and non-algebraic
scoped operations in Γ. First we construct its underlying endofunctor. When C
is Set, the intuition for programs PA is that they are terms inductively built
from the following inference rules:
a∈A
o∈Σn k: n→PA
s∈Γn p: n→PX k: X →PA
var(a) ∈PA
o(k) ∈PA
{s(p); k}∈PA
where n ranges over finite sets and o ∈Σn represents an algebraic operation
of |n|arguments, and similarly s ∈Γn is a scoped operation that creates |n|
scopes. The diﬀerence between algebraic and scoped operations is manifested by
an additional explicit continuation k in the third rule, as it is not the case that
sequentially composing s(p) with k equals s(p; k) like for algebraic operations,
so the continuation for scoped operations must be explicitly kept in the syntax.
When C is any lfp category, these rules translate to the following recursive
equation for the functor P : C →C:
PA∼
= A+ Σ(PA) +
X:C
C(X,PA)
where the existentially quantified X in the third rule is translated to a coend
X:C in C [44]. Moreover, the coend in (11) is isomorphic to Γ(P(PA)) because
by the coend formula of Kan extension, it exactly computes LanI(ΓP)(PA), i.e.
the left Kan-extension of ΓP along the identity functor I : C →C, and clearly
LanI(ΓP) = ΓP. Thus (11) is equivalent to
PA∼
= A+ Σ(PA) + Γ(P(PX)) (12)
which is exactly the Prog Σ Γ datatype that we saw in the Haskell implementa-
tion (7). To obtain a solution to (12), we construct a (higher-order) endofunctor
G : Endof(C) →Endof(C) to represent the Grammar where Endof(C) is the
category of finitary endofunctors on C:
G= Id + Σ◦−+ Γ◦−◦− (13)
where Id : C →C is the identity functor. Then Lemma 1 is applicable be-
cause Endof(C) has all small colimits since colimits in functor categories can be
computed pointwise and C has all small colimits. Furthermore, G preserves all
filtered colimits, in particular colimits of ω-chains, because −◦= : Endof(C) ×
Endof(C) →Endof(C) is finitary following from direct verification. Since initial
algebras are precisely free algebras generated by the initial object, by Lemma 1,
there is an initial G-algebra ⟨P : Endof(C), in : GP →P⟩and in is an isomor-
phism. Thus P obtained in this way is indeed a solution to (12)—the endofunctor
modelling the syntax of programs with algebraic and scoped operations.
14
Monadic Structure of P Next we equip the finitary endofunctor P with a monad
structure. This can be done in several ways, either by the general result about
Σ-monoids [17, 14] in Endof(C), or by [46, Theorem 4.3], or by the following
relatively straightforward argument in [49]: by the ‘diagonal rule’ of computing
initial algebras by Backhouse et al. [4], P= µG (13) is isomorphic to P′
=
µX. Id+Σ◦X+Γ◦P◦X. Note that P′is exactly (Σ+Γ◦P)∗as endofunctors
by Lemma 1, thus
P∼
= (Σ+ Γ◦P)∗
: Endof(C) (14)
Then we equip P with the same monad structure as the ordinary free monad
(Σ+ Γ◦P)∗. The implementation in (7) is exactly this monad structure.
3.3 Functorial Algebras of Scoped Operations
To interpret the monad P (12) modelling the syntax of scoped operations, it
is natural to expect that semantics is given by G-algebras on Endof(C) so that
interpretation is then the catamorphisms from µG to G-algebras. And follow-
ing the adjoint-theoretic approach (10), we would like to have an adjunction
⊣
G-Alg C
such that the induced monad is isomorphic to P. However, there
seems no natural way to construct such an adjunction unless we replace G-
algebras with a slight extension of it, which we referred to as functorial algebras,
as the notion for giving semantics to scoped operations. In the following, we first
define functorial algebras formally (Definition 1) and then show the adjunction
between the category of functorial algebras and the base category (Theorem 1),
which allows us to interpret P with functorial algebras.
A functorial algebra is carried by an endofunctor H : C →C with additionally
an object Xin C. The endofunctor Halso comes with a morphism αG : GH →H
in Endof(C), and the object X is equipped with a morphism αI : ΣX+ΓHX →
X in C. The intuition is that given a program of type PX∼
= X + Σ(PX) +
Γ(P(PX)), the middle P in ΓPP corresponds to the part of a program enclosed
by some scoped operations (i.e. the pin {s(p)>> =k}), and this part of the program
is interpreted by H with αG. After the enclosed part is interpreted, αI interprets
the outermost layer of the program by X with αI in the same way as interpreting
free monads of algebraic operations. More precisely, let I : Endof(C) ×C →C
be a bi-functor such that 3
IHX= ΣX+ Γ(HX) Iσf= Σf + Γ(σ◦f) (15)
for all H : Endof(C) and X : C and all morphisms σ: H →H′and f : X →X′
.
Then we define an endofunctor Fn : Endof(C) ×C →Endof(C) ×C such that
Fn⟨H, X⟩= ⟨GH, IHX⟩ (16)
3 The first argument H to I is written as subscript so that we have a more compact
notation I∗
H when taking the free monad of IH : CC with the first argument fixed.
15
f
Definition 1 (Functorial Algebras). A functorial algebra is an object ⟨H, X⟩
in Endof(C)×C paired with a structure map Fn⟨H, X⟩→⟨H, X⟩, or equivalently
it is a quadruple
H : Endof(C), X : C, αG : GH →H, αI : ΣX+ Γ(HX) →X
where GH= Id +Σ◦H+Γ◦H◦H. Morphisms between two functorial algebras
⟨H1, X1, αG
1 , αI
1⟩and ⟨H2, X2, αG
2 , αI
2⟩are pairs ⟨σ : H1 →H2, f : X1 →X2⟩
making the following diagrams commute:
αG
1
GH1 H1
Gσ σ
GH2 H2
αG
2
αI
1
ΣX1 + Γ(H1X1) X1
Σf+Γ(σ◦f)
ΣX2 + Γ(H2X2) X2
αI
2
Functorial algebras and their morphisms form a category Fn-Alg.
Example 8. We reformulate our programming example of nondeterministic choice
with once shown Example 4 in the formal definition. Let C= Set in this example
and 1 = {⋆}be some singleton set. We define signature endofunctors
ΣX = 1 + X×X ΓX= X
so that Σ represents nullary algebraic operation fail and binary algebraic opera-
tion or, and Γ represents the unary scoped operation once that creates one scope.
Let List : Set →Set be the endofunctor mapping a set X to the set of finite lists
with elements from X. We define natural transformations αΣ : Σ◦List →List
and αΓ : Γ ◦List ◦List →List by
αΣ
X(ι1 ⋆) = nil, αΣ
X(ι2 ⟨x, y⟩) = x++y, αΓ
X(nil) = nil, αΓ
X(cons x xs) = x
where nil is the empty list; ++ is list concatenation; and cons x xs is the list
with an element x in front of xs. Then for any set X, ⟨List, List X ⟩carries a
functorial algebra with structure maps
αG = [ηList
, αΣ
, αΓ] : GList →List αI = [αΣ
X, αΓ
X] : IList X →X (17)
where ηList : Id →List wraps any element into a singleton list.
The last example exhibits that one can define a functorial algebra carried
by ⟨H, HX⟩from a G-algebra on H : Endof(C) by simply choosing the object
component to be HX for an arbitrary X : C. In other words, there is a faithful
functor G-Alg →Fn-Alg, which results in functorial algebras that interpret the
outermost layer of a program—the part not enclosed by any scoped operation—
in the same way as the inner layers. But in general, the object component of
functorial algebras oﬀers the flexibility that the outermost layer can be inter-
preted diﬀerently from the inner layers, as in the following example.
Example 9. Continuing Example 8, if one is only interested in the final number
of possible outcomes, then one can define a functorial algebra ⟨List, N, αG
, αI⟩
where αG is (17) and αI(ι1 (ι1⋆)) = 0,
αI(ι1 (ι2⟨x, y⟩)) = x+ y, αI(ι2 nil) = 0, αI(ι2 (cons n ns)) = n
16
3.4 Interpreting with Functorial Algebras
In the rest of this section we show how functorial algebras can be used to in-
terpret programs PA (12) with scoped operations. We first construct a simple
adjunction ↑⊣↓between the base category C and Endof(C) ×C, which is then
composed with the free/forgetful adjunction FreeFn ⊣UFn between Endof(C)×C
and Fn-Alg for the functor Fn (16). The resulting adjunction (18) is proven to
induce a monad T isomorphic to P (Theorem 1), and by the adjoint-theoretic
approach to syntax and semantics (10), this adjunction provides a means to
interpret scoped operations modelled with the monad P (Theorem 2).
First we define functor ↑: C →Endof(C) ×C such that ↑X= ⟨0, X⟩where
0 : Endof(C) is the initial endofunctor—the constant functor sending everything
to the initial object in C. The functor ↑is left adjoint to the projection functor
↓: Endof(C) ×C →C of the second component.
Then we would like to compose ↑ ⊣ ↓with the free-forgetful adjunction
FreeFn ⊣UFn for the endofunctor Fn (16) on Endof(C) ×C, and the latter ad-
junction indeed exists.
Lemma 2. The endofunctor Fn (16) on Endof(C) ×C has free algebras, i.e.
there is a functor FreeFn : Endof(C) ×C →Fn-Alg left adjoint to the forgetful
functor UFn : Fn-Alg →Endof(C) ×C.
Proof sketch. It can be verified that Fn is finitary and then we apply Lemma 1.
A detailed proof can be found in Appendix C.
These two adjunctions are depicted in the following diagram:
FreeFn
↑
⊥
⊥
Fn-Alg Endof(C) ×C C
UFn
↓
T (18)
and we compose them to obtain an adjunction FreeFn ↑⊣↓UFn between Fn-Alg
and C, giving rise to a monad T= ↓UFnFreeFn ↑. In the rest of this section, we
prove that T is isomorphic to P (11) in the category of monads, which is crucial
in this paper, since it allows us to interpret scoped operations modelled by the
monad P with functorial algebras Fn-Alg.
We first establish a technical lemma characterising the free Fn-algebra on the
product category Endof(C) ×C in terms of the free algebras in C and Endof(C).
Lemma 3. There is a natural isomorphism between FreeFn and the following
FreeFn⟨H, X⟩= G∗H : Endof(C), (IG∗H)∗X : C, opG∗
H , op(IG∗H)∗
X
where opG∗
H : G(G∗H) →G∗H and op(IG∗H)∗
X : IG∗H((IG∗H)∗X) →(IG∗H)∗X
are the structure maps of the free G-algebra and IG∗H-algebra respectively.
Proof sketch. It follows from the formula of computing free algebras from initial
algebras Lemma 1. A detailed proof can be found in the appendix.
17
Theorem 1. Monads P (12) and T (18) are isomorphic as monads.
Proof sketch. By (14) and (15), P∼
= (Σ+ΓP)∗= (IP)∗as monads, it is suﬃcient
to show that T∼
= (IP)∗ as monads, which follows from Lemma 3 by careful
calculation. A detailed proof can be found in Appendix C.
Remark 3. In general, the adjunction FreeFn ↑⊣↓UFn is not monadic since the
right adjoint ↓UFn does not reflect isomorphisms, which is a necessary condi-
tion for it to be monadic by Beck’s monadicity theorem [44]. This entails that
the category Fn-Alg of functorial algebras is not equivalent to the category of
Eilenberg-Moore algebras. Nonetheless, as we will see later in Section 4, functo-
rial algebras and Eilenberg-Moore algebras have the same expressive power for
interpreting scoped operations in the base category.
The isomorphism established Theorem 1 enables us to interpret programs
modelled by the monad P using functorial algebras following (10): for any func-
torial algebra ⟨H, X, αG
, αI⟩(Definition 1), and any morphism g : A →X in
the base category C, there is a morphism
handle⟨H,X,αG
,αI⟩g= ↓UFn(ϵ⟨H,X,αG
,αI⟩·FreeFn ↑g) : TA∼
= PA→X (19)
which interprets programs PA with the functorial algebra ⟨H, X, αG
, αI⟩. Fur-
thermore, we can derive the following recursive formula (20) for this interpreta-
tion morphism, which is exactly the Haskell implementation in Figure 1.
Theorem 2 (Interpreting with Functorial Algebras). For any functorial
algebra α = ⟨H, X, αG
, αI⟩as in Definition 1, and any morphism g : A →X
for some A in the base category C, let h= αG : P →H be the catamorphism
from the initial G-algebra P to the G-algebra αG : GH →H. The interpretation
of PA with this algebra α and g satisfies
handleα g= [g, αI
Σ·Σ(handleα g), αI
Γ·ΓhX·ΓP(handleα g)]·in◦
A (20)
where in◦
: P →Id + Σ◦P + Γ◦P◦P is the isomorphism between P and GP;
morphisms αI
Σ = αI
·ι1 : ΣX →X and αI
Γ = αI
·ι2 : ΓHX →X are the two
components of αI : ΣX+ ΓHX →X.
Proof sketch. It can be calculated by plugging the formula of ϵfor the adjunction
↓UFn ⊣FreeFn ↑(Lemma 1) into (19).
To summarise, we have defined a notion of functorial algebras that we use to
handle scoped operations. The heart of the development is the adjunction (18)
that induces a monad isomorphic to the monad P (12) that models the syntax of
programs with scoped operations, following which we derive a recursive formula
(20) that interprets programs with functor algebras. The formula is exactly the
implementation in Figure 1: the datatype EndoAlg represents the αG in (20);
datatype BaseAlg corresponds to αI; function hcata implements αG
.
18
4 Comparing the Models of Scoped Operations
Functorial algebras are not the only option for interpreting scoped operations. In
this section we compare functorial algebras with two other approaches, one being
Pir´og et al. [49]’s indexed algebras and the other one being Eilenberg-Moore (EM)
algebras of the monad P (12), which simulate scoped operations with algebraic
operations. After a brief description of these two kinds of algebras, we compare
them and show that their expressive power is in fact equivalent.
4.1 Interpreting Scoped Operations with Eilenberg-Moore Algebras
In standard algebraic eﬀects, handlers are just Σ-algebras for some signature
functor Σ : C →C, and it is well known that the category Σ-Alg of Σ-algebras
is equivalent to the category CΣ∗of EM algebras of the monad Σ∗. Thus handlers
of algebraic operations are exactly EM algebras of the monad Σ∗modelling the
syntax of algebraic operations. This observation suggests that we may also use
EM algebras of the monad P (12) as the notion of handlers for scoped operations.
Lemma 4. EM algebras of P are equivalent to (Σ+ Γ◦P)-algebras. In other
words, an EM algebra of P is equivalently a tuple
⟨X : C, αΣ : ΣX →X, αΓ : Γ(PX) →X⟩ (21)
Proof. Recall that the monad structure of P (14) is exactly the monad structure
of the free monad (Σ+ Γ◦P)∗, and therefore they have the same EM algebras.
Moreover, EM algebras of (Σ + Γ◦P)∗ are equivalent to plain (Σ + Γ◦P)-
algebras by the monadicity of the free-forgetful adjunction.
Thus we obtain a way of interpreting scoped operations based on the free-
forgetful adjunction FreeΣ+Γ◦P ⊣UΣ+Γ◦P: given an EM algebra α= ⟨X, αΣ, αΓ⟩
of P as in (21), then for any A: C and morphism g: A→X, the interpretation
of PA by g and this EM algebra is
handleα g= UΣ+Γ◦P(ϵα·FreeΣ+Γ◦P g) : PA∼
= (Σ+ Γ◦P)∗A→X (22)
The formula (22) can also be turned into a recursive form:
handleα g= [g, αΣ·Σ(handleα g), αΓ·ΓP(handleα g)]·in◦
A (23)
that suits implementation (see Figure 3 and Appendix D for more details).
Interpreting scoped operation with EM algebras can be understood as sim-
ulating scoped operations with algebraic operations and general recursion: a
signature (Σ, Γ) of algebraic-and-scoped operations is simulated by a signature
(Σ+Γ◦P) of algebraic operations where P is recursively given by (Σ+Γ◦P)∗. In
this way, one can simulate scoped operation in languages implementing algebraic
eﬀects that allow signatures of operation to be recursive, such as [21, 7, 39], but
not the original design by Plotkin and Pretnar [55], which requires signatures of
operations to mention only some base types.
19
The downside of this simulating approach is that the denotational seman-
tics of the language becomes more complex and usually involves solving some
domain-theoretic recursive equations, like in [7]. Moreover, this approach typi-
cally requires handlers to be defined with general recursion, which obscures the
inherent structure of scoped operations, making reasoning about handlers of
scoped operations more diﬃcult.
4.2 Indexed Algebras of Scoped Eﬀects
Indexed algebras of scoped operations by Pir´og et al. [49] are yet another way
of interpreting scoped operations. They are based on the following adjunction:
FreeIx
⊣
⊣
Ix-Alg C|N| C
UIx
(24)
where C|N| is the functor category from the discrete category |N|of natural
numbers to the base category C. That is to say, an object in C|N|is a family of
objects Ai in C indexed by natural numbers i∈|N|, and a morphism τ : A→B
in C|N|is a family of morphisms τi : Ai →Bi in C with no coherence conditions
between the levels. An endofunctor Ix : C|N| →C|N| is defined to characterise
indexed algebras:
IxA= Σ◦A+ Γ◦(▹A) + (◃A)
where ▹ and ◃ are functors C|N|→C|N|shifting indices such that (▹A)i = Ai+1
and (◃A)0 = 0 and (◃A)i+1 = Ai. Then objects in Ix-Alg are called indexed
algebras. Furthermore, since a morphism (◃A) →A is in bijection with A →
(▹A), an indexed algebra can be given by the following tuple:
⟨A: C|N|
, a: Σ◦A→A, d: Γ(▹A) →A, p: A→▹A⟩ (25)
The operational intuition for it is that the carrier Ai at level i interprets the
part of syntax enclosed by i layers of scopes, and when interpreting a scoped
operation Γ(P(PX)) at layer i, the part of syntax outside the scope is first
interpreted, resulting in Γ(PAi), and then the indexed algebra provides a way
p to promote the carrier to the next level, resulting in Γ(PAi+1). After the
inner layer is also interpreted as ΓAi+1, the indexed algebra provides a way dto
demote the carrier, producing Ai again. Additionally the morphism ainterprets
ordinary algebraic operations.
Example 10. Example 8 for nondeterministic choice with once can be re-expressed
with an indexed algebra as follows. For any set X, we define an indexed object
A: C|N| by A0 = List X and Ai+1 = List Ai. The object A carries an indexed
algebra with the following structure maps: for all i∈N, ai(ι1 ⋆) = nil and
ai(ι2 ⟨x, y⟩) = x++y, di(nil) = nil, di(cons x xs) = x, pi(x) = cons x nil
20
The adjunction FreeIx ⊣UIx in (24) is the free-forgetful adjunction for Ix
on C|N|. The other adjunction ⊣ is given by A= A0, ( X)0 = X, and
( X)i+1 = 0 for all i ∈N. Importantly, Pir´og et al. [49] show that the monad
induced by the adjunction (24) is isomorphic to monad P (12), thus indexed
algebras can also be used to interpret scoped operations
handle⟨A,a,d,p⟩g= UIx(ϵ⟨A,a,d,p⟩·FreeIx g) (26)
in the same way as what we do for functorial algebras in Section 3.4. Interpreting
with indexed algebras can also be implemented in Haskell with GHC’s DataKinds
extension for type-level natural numbers (which can be found in Figure 4).
4.3 Comparison of Resolutions
Now we come back to the real subject of this section—comparing the expressivity
of the three ways for interpreting scoped operations. Specifically, we construct
comparison functors between the respective categories of the three kinds of al-
gebras, which translate one kind of algebras to another in a way preserving the
induced interpretation in the base category. Categorically, the three kinds of
algebras correspond to three resolutions of the monad P, which form a cate-
gory of resolutions (Definition 2) with comparison functors as morphisms. In
this category, the Eilenberg-Moore resolution is the terminal object, and thus
it automatically gives us comparison functors translating other kinds of alge-
bras to EM algebras. To complete the circle of translations, we then construct
comparison functors KEM
Fn : CP →Fn-Alg translating EM algebras to functorial
ones (Section 4.4) and KFn
Ix : Fn-Alg →Ix-Alg translating functorial algebras to
indexed ones (Section 4.5). The situation is pictured in Figure 2.
Definition 2 (Resolutions and Comparison Functors [38]). Given a
monad M on C, the category Res(M) of resolutions of M has as objects ad-
junctions ⟨D, L ⊣R : D →C, η, ϵ⟩such that the induced monad RL is M. A
morphism from a resolution ⟨D, L ⊣R, η, ϵ⟩to ⟨D′
, L′ ⊣R′
, η′
, ϵ′⟩is a functor
K : D →D′, called a comparison functor, such that it commutes with the left
and right adjoints, i.e. KL= L′ and R′K= R.
We have seen adjunctions for indexed algebras, EM algebras and functorial
algebras respectively, each inducing the monad P up to isomorphism, so each of
them can be identified with an object in the category Res(E). For each resolution
⟨D, L, R, η, ϵ⟩, we have been using the objects D in D to interpret scoped opera-
tions modelled by P: for any morphism g : A→RD in C, the interpretation of
PA by D and g is handleD g= R(ϵD·Lg) : PA= RLA →RD. Crucially, we
show that interpretations are preserved by comparison functors.
Lemma 5 (Preservation of Interpretation). Let K : D →D′ be any com-
parison functor between resolutions ⟨D, L, R, η, ϵ⟩and ⟨D′
, L′
, R′
, η′
, ϵ′⟩of some
monad M : C →C. For any object D in D and any g: A→RD in C,
handleD g= handleKD g: MA→RD(= R′KD) (27)
where each side interprets MA using L⊣R and L′⊣R′ respectively.
21
Proof. Since L⊣Rand L′⊣R′induce the same monad, their unit must coincide
η= η′. Together with the commutativity properties KL= L′ and R′K= R,
it makes a comparison functor a special case of a map of adjunctions. Then by
Proposition 1 in [44, page 99], it holds that Kϵ= ϵ′K, and we have
handleKD g= R′(ϵ
′
KD·L′g) = R′(KϵD·L′g)
= RϵD·R′L′g= RϵD·RLg= handleD g
This lemma implies that if there is a comparison functor K from some reso-
lution L⊣R: D →C to L′⊣R′: D′→C of the monad P, then K can translate
a D object to a D′ object that preserves the induced interpretation. Thus the
expressive power of D for interpreting P is not greater than D′, in the sense
that every handleD g that one can obtain from D in D can also be obtained by
an algebra KD in D′. Thus the three kinds of algebras for interpreting scoped
operations have the same expressivity if we can construct a circle of comparison
functors between their categories, which is what we do in the following.
Translating to EM Algebras As shown in [44], an important property of the
Eilenberg-Moore adjunction is that it is the terminal object in the category
Res(M) for any monad M, which means that there uniquely exists a comparison
functor from every resolution to the Eilenberg-Moore resolution. Specifically,
given a resolution ⟨D, L, R, η, ϵ⟩of a monad M, the unique comparison functor
K from D to the category CM of the Eilenberg-Moore algebras is
KD= M(RD) = RLRD RϵD −−−→RD and K(D f −→D′) = Rf
Lemma 6. There uniquely exist comparison functors KIx
EM : Ix-Alg →CP and
KFn
EM : Fn-Alg → CP from the resolutions of indexed algebras and functorial
algebras to the resolution of EM algebras.
4.4 Translating EM Algebras to Functorial Algebras
Now we construct a comparison functor KEM
Fn : CP →Fn-Alg translating EM
algebras to functorial ones. The idea is straightforward: given an EM algebra X,
we map it to the functorial algebra with X for interpreting the outermost layer
and the functor P for interpreting the inner layers, which essentially leaves the
inner layers uninterpreted before they get to the outermost layer.
Since CP is isomorphic to (Σ+Γ◦P)-Alg, we can define KEM
Fn on (Σ+Γ◦P)-
algebras instead. Given any ⟨X : C, α : (Σ+ Γ◦P)X →X⟩, it is mapped by
KEM
Fn to the functorial algebra
⟨P, X, in : GP →P, α: (Σ+ Γ◦P)X →X⟩
and for any morphism f in (Σ+ Γ◦P)-Alg, it is mapped to ⟨idP, f⟩. To show
KEM
Fn is a comparison functor, we only need to show that it commutes with the
left and right adjoints of both resolutions. Details can be found in Appendix C.
Lemma 7. Functor KEM
Fn is a comparison functor from the Eilenberg-Moore res-
olution of P to the resolution FreeFn ↑⊣↓UFn of functorial algebras.
22
4.5 Translating Functorial Algebras to Indexed Algebras
KIx
At this point we have comparison functors Ix-Alg
EM −−→CP KEM
Fn −−→Fn-Alg. To
complete the circle of translations, we construct a comparison functor KFn
Ix :
Fn-Alg →Ix-Alg in this subsection. The idea of this translation is that given
a functorial algebra carried by endofunctor H : CC and object X : C, we map
it to an indexed algebra by iterating the endofunctor H on X. More precisely,
KFn
Ix : Fn-Alg →Ix-Alg maps a functorial algebra
⟨H : CC
, X : C, αG : Id + Σ◦H+ Γ◦H◦H →H, αI : ΣX+ ΓHX →X⟩
to an indexed algebra carried by A : C|N| such that Ai = HiX, i.e. iterating
H i-times on X. The structure maps of this indexed algebra ⟨a: ΣA→A, d:
Γ(▹A) →A, p: A→(▹A)⟩are given by
a0 = (αI
·ι1) : ΣX →X ai+1 = (αG
HiX·ι2) : ΣHHiX →Hi+1X
d0 = (αI
·ι2) : ΓHX →X di+1 = (αG
HiX·ι3) : ΓHHHiX →Hi+1X
and pi = αG
HiX·ι1 : HiX →HHiX. On morphisms, KFn
Ix maps a morphism
⟨τ : H →H′
, f : X →X′⟩in Fn-Alg to σ : HiX →H′iX′in Ix-Alg such that
σ0 = f and σi+1 = τ◦σi where ◦is horizontal composition.
Lemma 8. KFn
Ix is a comparison functor from the resolution FreeFn ↑⊣↓UFn of
functorial algebras to the resolution FreeIx ⊣ UIx of indexed algebras.
Proof sketch. Again we only need to show the required commutativities for KFn
Ix
to be a comparison functor: ↓UFn
∼
= UIxKFn
Ix and KFn
Ix FreeFn ↑∼
= FreeIx . The
first one is easy, and the second one follows from Pir´og et al. [49]’s explicit
characterisation of FreeIx X. More details can be found in Appendix C.
Since comparison functors preserve interpretation (Lemma 5), the lemma
above implies that the expressivity of functorial algebras is not greater than
indexed ones. Together with the comparison functors defined earlier, we con-
clude that the three kinds of algebras—indexed, functorial and Eilenberg-Moore
algebras—have the same expressivity for interpreting scoped operations. Figure 2
summarises the comparison functors and resolutions that we have studied.
Remark 4. Although the three kinds of algebras have the same expressivity in
theory, they structure the interpretation of scoped operations in diﬀerent ways:
EM algebras impose no constraint on how the part of syntax enclosed by scopes is
handled; indexed algebras demand them to be handled layer by layer but impose
no coherent conditions between the layers; functorial algebras additionally force
all inner layers must be handled in a uniform way by an endofunctor.
On the whole, it is a trade-oﬀ simplicity and structuredness: EM algebras
are the simplest for implementation, whereas the structuredness of functorial
algebras make them easier to reason about. This is another instance of the
preference for structured programming over unstructured language features, in
the same way as structured loops being favoured over goto, although they have
the same expressivity in theory [13].
23
5 Fusion Laws of Interpretation
A crucial advantage of the adjoint-theoretic approach to syntax and semantics is
that the naturality of an adjunction directly oﬀers fusion laws of interpretation
that fuse a morphism after an interpretation into a single interpretation, which
have proven to be a powerful tool for reasoning and optimisation [69, 68, 12, 24,
72, 74]. In this section, we present the fusion law for functorial algebras.
5.1 Fusion Laws of Interpretation
Recall that given any resolution L⊣Rwith counit ϵof some monad M : C →C
where L: C →D, for any g: A→RD, we have an interpretation morphism
handleD g= R(ϵD·Lg) : MA→RD
Then whenever we have a morphism in the form of (f·handleD g)—an interpre-
tation followed by some morphism—the following fusion law allows one to fuse
it into a single interpretation morphism.
Lemma 9 (Interpretation Fusion). Assume L⊣Ris a resolution of monad
M : C →C where L: C →D. For every D : D, g : A→RD and f : RD→X,
if there is some D′ and h: D→D′ in D such that RD′
= X and Rh= f, then
f·handleD g= handleD′ (f·g) (28)
Proof. We have f·handleD g= Rh·R(ϵD·Lg) = R(h·ϵD·Lg). Then by the
naturality of the counit ϵ, R(h·ϵD·Lg) = R(ϵD′·L(Rh·g)) = handleD′ (f·g).
Applying the lemma to the three resolutions of P gives us three fusion laws:
for any D: D where D ∈{Ix-Alg, Fn-Alg, CP}, one can fuse f·handleD g into
a single interpretation if one can make f a D-homomorphism. Particularly, the
following is the fusion law for functorial algebras.
Corollary (Fusion Law for Functorial Algebras). Letˆ
α1 = ⟨H, X1, αG
1 , αI
2⟩
be a functorial algebra (Definition 1) and g : A →X1, f : X1 →X2 be any
morphisms in C. If there is a functorial algebraˆ
α2 = ⟨H2, X2, αG
2 , αI
2⟩and a
functorial algebra morphism ⟨σ: H1 →H2, h: X1 →X2⟩, then
f·handleˆ
α1 g= handleˆ
α2 (f·g)
Example 11. Let ˆ α be the functorial algebra of nondeterminism with once in
Example 8 and len : List A →N be the function mapping a list to its length.
Then using the fusion law, len·handleˆ
α g= handleˆ
β (len·g) if we can find a
ˆ
suitable functorial algebra
β : Fn-Alg and h: ˆ α→ˆ
β s.t. ↓UFnh= len. In fact, a
ˆ
suitable
β is just the functorial algebra in Example 9 and h= ⟨id, len⟩.
24
Example 12. Although Pir´og et al. [49] propose the adjunction (24) to interpret
scoped operations with indexed algebras, their Haskell implementation is not
a faithful implementation of the interpretation morphism (26), but rather a
more eﬃcient one skipping the step of transforming P to the isomorphic free
indexed algebra ( UIxFreeIx ). However, it is previously unclear whether this
implementation indeed coincides with the interpretation morphism (26) due to
the discrepancy between the syntax monad P and indexed algebras.
This issue is in fact one of the original motivations for us to develop functo-
rial algebras—a way to interpret P that directly follows the syntactic structure.
Using the comparison functors to transform between indexed and functorial al-
gebras, we can reason about Pir´og et al. [49]’s implementation with functorial
algebras, and its correctness can be established using fusion laws. Due to space
limit, we do not expand on this and refer interested readers to Appendix E.
6 Related Work
The most closely related work is of course that of Pir´og et al. [49] on categor-
ical models of scoped eﬀects. That work in turn builds on Wu et al. [73] who
introduced the notion of scoped eﬀects after identifying modularity problems
with using algebraic eﬀect handlers for catching exceptions [55]. Scoped eﬀects
have found their way into several Haskell implementations of algebraic eﬀects
and handlers [59, 45, 35].
Eﬀect Handlers and Modularity Spivey [63], Moggi [47] and Wadler [70] initiated
monads for modeling and programming with computational eﬀects. Soon after,
the desire arose to define complex monads by combining modular definitions
of individual eﬀects [66, 29], and monad transformers were developed to meet
this need [42]. Yet, several years later, algebraic eﬀects were proposed as an
alternative more structured approach for defining and combining computational
eﬀects [51, 52, 25]. The addition of handlers [55] has made them practical for
implementation and many languages and libraries have been developed since.
Schrijvers et al. [61] have characterized modular handlers by means of modular
carriers, and shown that they correspond to a subclass of monad transformers.
Forster et al. [19] have also shown that algebraic eﬀects, monads and delimited
control are macro-expressible in terms of each other in an untyped language but
not in a simply typed language.
Scoped operations are generally not algebraic operations in the original design
of algebraic eﬀects [51], but as we have seen in Section 4.1, an alternative view
on Eilenberg-Moore algebras of scoped operations is regarding them as handlers
of algebraic operations of signature Σ + ΓP. However, the functor Σ + ΓP
involves the type P modelling computations, and thus it is not a valid signature
of algebraic eﬀects in the original design of eﬀect handlers [54, 55], in which the
signature of algebraic eﬀects can only be built from some base types to avoid
the interdependence of the denotations of signature functors and computations.
In spite of that, many later implementations of eﬀect handlers such as Eff [7],
25
Koka [39] and Frank [43] do not impose this restriction on signature functors
(at the cost that the denotational semantics involves solving recursive domain-
theoretic equations), and thus scoped operations can be implemented in these
languages with EM algebras as handlers.
Other variations of scoped eﬀects have recently been suggested. Recently,
Poulsen et al. [56] and van den Berg et al. [9] have proposed a notion of staged
or latent eﬀect, which is a variant of scoped eﬀects, for modelling the deferred
execution of computations inside lambda abstractions and similar constructs.
Ahman and Pretnar [3] investigate asynchronous eﬀects, and they note that
interrupt handlers are in fact scoped operations. We have not yet investigated
this in our framework, but it will be an interesting use case.
Abstract Syntax This work focusses on the problem of abstract syntax and se-
mantics of programs. The practical benefit of abstract syntax is that it allows for
generic programming in languages like Haskell that have support for, e.g. type
classes, gadts [28] and so on. As an example, Swierstra [67] showed that it is
possible to modularly create compilers by formalising syntax in Haskell.
The problem of formalising abstract syntax categorically for operations with
variable binding was first addressed by Fiore et al. [17, 18]. Subsequently, Ghani
et al. [20] model the abstract syntax of explicit substitutions as an initial algebra
in the endofunctor category and show that it is a monad. Pir´og et al. [49] and
this paper use a monad P, which is a slight generalisation of the monad of
explicit substitutions, to model the syntax of scoped operations. The datatype
underlying P is an instance of nested datatypes studied by Bird and Paterson
[10] and Johann and Ghani [27].
In this paper we have not treated equations on eﬀectful operations, which
are both theoretically and practically important. Plotkin and Power [51] show
that theories of various eﬀects with suitable equations determine their corre-
sponding monads, and later Hyland et al. [25] show that certain combinations
of eﬀect theories are equivalent to monad transformers. Equations are also used
for reasoning about programs with algebraic eﬀects and handlers [53, 74, 37].
Possible ways to extend scoped eﬀects with equations include the approach in
[33] (Remark 1), the categorical framework of equational systems [15], second
order Lawvere theories [16], and syntactic frameworks like [65].
7 Conclusion
The motivation of this work is to develop a modular approach to the syntax
and semantics of scoped operations. We believe our proposal, functorial alge-
bras, is at a sweet spot in the trade-oﬀ between structuredness and simplicity,
allowing practical examples of scoped operations to be programmed and rea-
soned about naturally, and implementable in modern functional languages such
as Haskell and OCaml. We put our model of interpreting scoped operations and
two existing models in the same framework of resolutions of the monad mod-
elling syntax, and by constructing interpretation-preserving functors between
26
the three kinds of algebras, we showed that they have equivalent expressivity
for interpreting scoped operation, although they form non-equivalent categories.
The uniform theoretical framework also gave rise to fusion laws of interpretation
in a straightforward way.
There are two strains of work that should be pursued from here. The first one
would be investigating modularity of algebras of scoped operations, in particular,
ways to forward scoped operations that is not handled. The second one would
be the design of a language supporting handlers of scoped operations natively
and its type system and operational semantics.
Bibliography
[1] Adamek, J., Rosicky, J.: Locally Presentable and Accessible Categories.
London Mathematical Society Lecture Note Series, Cambridge University
Press (1994), https://doi.org/10.1017/CBO9780511600579
[2] Ad´amek, J.: Free algebras and automata realizations in the language of
categories. Commentationes Mathematicae Universitatis Carolinae 015(4),
589–602 (1974), URL http://eudml.org/doc/16649
[3] Ahman, D., Pretnar, M.: Asynchronous eﬀects. Proc. ACM Program. Lang.
(POPL) (2021), https://doi.org/10.1145/3434305, URL https://doi.org/
10.1145/3434305
[4] Backhouse, R., Bijsterveld, M., van Geldrop, R., van der Woude, J.: Cat-
egorical fixed point calculus. In: Pitt, D., Rydeheard, D.E., Johnstone,
P. (eds.) Category Theory and Computer Science, pp. 159–179, Springer
Berlin Heidelberg, Berlin, Heidelberg (1995), ISBN 978-3-540-44661-3,
https://doi.org/10.1007/3-540-60164-3 25
[5] Barr, M.: Coequalizers and free triples. Mathematische Zeitschrift (1970)
[6] Barr, M., Wells, C.: Category theory for computing science, vol. 1. Prentice
Hall New York (1990)
[7] Bauer, A., Pretnar, M.: An eﬀect system for algebraic eﬀects and handlers.
Logical Methods in Computer Science 10(4) (Dec 2014), ISSN 1860-5974,
https://doi.org/10.2168/lmcs-10(4:9)2014
[8] Bauer, A., Pretnar, M.: Programming with algebraic eﬀects and han-
dlers. J. Log. Algebraic Methods Program. 84(1), 108–123 (2015),
https://doi.org/10.1016/j.jlamp.2014.02.001, URL https://doi.org/10.
1016/j.jlamp.2014.02.001
[9] van den Berg, B., Schrijvers, T., Bach-Poulsen, C., Wu, N.: Latent eﬀects
for reusable language components: Extended version (2021)
[10] Bird, R.S., Paterson, R.: Generalised folds for nested datatypes. For-
mal Aspects Comput. (1999), https://doi.org/10.1007/s001650050047, URL
https://doi.org/10.1007/s001650050047
[11] CLAESSEN, K.: A poor man’s concurrency monad. Jour-
nal of Functional Programming 9(3), 313–323 (1999),
https://doi.org/10.1017/S0956796899003342
[12] Coutts, D., Leshchinskiy, R., Stewart, D.: Stream fusion: From lists to
streams to nothing at all. SIGPLAN Not. 42(9), 315–326 (Oct 2007),
ISSN 0362-1340, https://doi.org/10.1145/1291220.1291199, URL https:
//doi.org/10.1145/1291220.1291199
[13] Dijkstra, E.W.: Letters to the editor: Go to statement considered
harmful. Commun. ACM 11(3), 147–148 (Mar 1968), ISSN 0001-0782,
https://doi.org/10.1145/362929.362947
[14] Fiore, M., Hur, C.K.: On the construction of free algebras for equational sys-
tems. Theoretical Computer Science 410(18), 1704–1729 (2009), ISSN 0304-
28
3975, https://doi.org/https://doi.org/10.1016/j.tcs.2008.12.052, automata,
Languages and Programming (ICALP 2007)
[15] Fiore, M., Hur, C.K.: On the construction of free algebras for equational sys-
tems. Theoretical Computer Science 410(18), 1704–1729 (2009), ISSN 0304-
3975, https://doi.org/https://doi.org/10.1016/j.tcs.2008.12.052, automata,
Languages and Programming (ICALP 2007)
[16] Fiore, M., Mahmoud, O.: Second-order algebraic theories. In: Hlinˇen´y, P.,
Kuˇcera, A. (eds.) Mathematical Foundations of Computer Science 2010,
pp. 368–380, Springer Berlin Heidelberg, Berlin, Heidelberg (2010), ISBN
978-3-642-15155-2
[17] Fiore, M.P., Plotkin, G.D., Turi, D.: Abstract syntax and variable binding.
In: 14th Annual IEEE Symposium on Logic in Computer Science, Trento,
Italy, July 2-5, 1999 (1999)
[18] Fiore, M.P., Turi, D.: Semantics of name and value passing. In: 16th Annual
IEEE Symposium on Logic in Computer Science, Boston, Massachusetts,
USA, June 16-19, 2001, Proceedings (2001)
[19] Forster, Y., Kammar, O., Lindley, S., Pretnar, M.: On the ex-
pressive power of user-defined eﬀects: Eﬀect handlers, monadic re-
flection, delimited control. J. Funct. Program. 29, e15 (2019),
https://doi.org/10.1017/S0956796819000121, URL https://doi.org/10.
1017/S0956796819000121
[20] Ghani, N., Uustalu, T., Hamana, M.: Explicit substitutions and higher-order
syntax. High. Order Symb. Comput. (2006)
[21] Hillerstr¨om, D., Lindley, S.: Shallow Eﬀect Handlers. Lecture Notes
in Computer Science 11275 LNCS, 415–435 (2018), ISSN 16113349,
https://doi.org/10.1007/978-3-030-02768-1 22
[22] Hinze, R.: Prological features in a functional setting — axioms and imple-
mentations. In: Sato, M., Toyama, Y. (eds.) Proceedings of the Third Fuji
International Symposium on Functional and Logic Programming (FLOPS
’98), pp. 98–122, World Scientific, Singapore, New Jersey, London, Hong
Kong (apr 1998), ISBN ISBN 981-02-3384-1
[23] Hinze, R.: Adjoint folds and unfolds—an extended study. Science of Com-
puter Programming (2013)
[24] Hinze, R., Harper, T., James, D.W.H.: Theory and practice of fusion. In:
Hage, J., Moraz´an, M.T. (eds.) Implementation and Application of Func-
tional Languages, pp. 19–37, Springer Berlin Heidelberg, Berlin, Heidelberg
(2011), ISBN 978-3-642-24276-2, https://doi.org/10.1007/978-3-642-24276-
2 2
[25] Hyland, M., Plotkin, G., Power, J.: Combining eﬀects: Sum and ten-
sor. Theor. Comput. Sci. 357(1), 70–99 (Jul 2006), ISSN 0304-3975,
https://doi.org/10.1016/j.tcs.2006.03.013
[26] Hyland, M., Power, J.: The category theoretic understanding of uni-
versal algebra: Lawvere theories and monads. Electronic Notes in
Theoretical Computer Science 172, 437–458 (2007), ISSN 1571-0661,
https://doi.org/https://doi.org/10.1016/j.entcs.2007.02.019, computation,
Meaning, and Logic: Articles dedicated to Gordon Plotkin
29
[27] Johann, P., Ghani, N.: Initial algebra semantics is enough! In: Typed
Lambda Calculi and Applications, TLCA, Lecture Notes in Computer Sci-
ence, Springer (2007), https://doi.org/10.1007/978-3-540-73228-0 16, URL
https://doi.org/10.1007/978-3-540-73228-0_16
[28] Johann, P., Ghani, N.: Foundations for structured programming with
gadts. In: Necula, G.C., Wadler, P. (eds.) Proceedings of the 35th ACM
SIGPLAN-SIGACT Symposium on Principles of Programming Languages,
POPL 2008, San Francisco, California, USA, January 7-12, 2008, pp. 297–
308, ACM (2008), https://doi.org/10.1145/1328438.1328475, URL https:
//doi.org/10.1145/1328438.1328475
[29] Jones, M.P., Duponcheel, L.: Composing monads. Research Report
YALEU/DCS/RR-1004, Yale University, New Haven, Connecticut, USA
(December 1993), URL http://web.cecs.pdx.edu/~mpj/pubs/RR-1004.
pdf
[30] Kammar, O., Lindley, S., Oury, N.: Handlers in action. SIG-
PLAN Not. 48(9), 145–158 (Sep 2013), ISSN 0362-1340,
https://doi.org/10.1145/2544174.2500590, URL https://doi.org/
10.1145/2544174.2500590
[31] Kammar, O., Plotkin, G.D.: Algebraic foundations for eﬀect-dependent op-
timisations. In: Proceedings of the 39th Annual ACM SIGPLAN-SIGACT
Symposium on Principles of Programming Languages, p. 349–360, POPL
’12, Association for Computing Machinery, New York, NY, USA (2012),
ISBN 9781450310833, https://doi.org/10.1145/2103656.2103698
[32] Kelly, G., Power, A.: Adjunctions whose counits are coequaliz-
ers, and presentations of finitary enriched monads. Journal of
Pure and Applied Algebra 89(1), 163–179 (1993), ISSN 0022-4049,
https://doi.org/https://doi.org/10.1016/0022-4049(93)90092-8
[33] Kelly, G.M.: Structures defined by finite limits in the enriched context, i.
Cahiers de Topologie et G´eom´etrie Diﬀ´erentielle Cat´egoriques 23(1), 3–42
(1982), URL http://www.numdam.org/item/CTGDC_1982__23_1_3_0/
[34] Kidney, D.O., Wu, N.: Algebras for weighted search. Proc. ACM Program.
Lang. 5(ICFP) (Aug 2021), https://doi.org/10.1145/3473577, URL https:
//doi.org/10.1145/3473577
[35] King, A.: eﬀ – screaming fast extensible eﬀects for less (2019), https://
github.com/hasura/eff
[36] Kiselyov, O., Ishii, H.: Freer monads, more extensible eﬀects. p. 94–105,
Association for Computing Machinery, New York, NY, USA (Aug 2015),
ISBN 9781450338080, https://doi.org/10.1145/2804302.2804319
[37] KISELYOV, O., MU, S.C., SABRY, A.: Not by equations alone: Reasoning
with extensible eﬀects. Journal of Functional Programming 31, e2 (2021),
https://doi.org/10.1017/S0956796820000271
[38] Lambek, J., Scott, P.J.: Introduction to Higher Order Categorical Logic
(1986)
[39] Leijen, D.: Type directed compilation of row-typed algebraic eﬀects.
In: Proceedings of the 44th ACM SIGPLAN Symposium on Prin-
ciples of Programming Languages, p. 486–499, POPL 2017, Associa-
30
tion for Computing Machinery, New York, NY, USA (2017), ISBN
9781450346603, https://doi.org/10.1145/3009837.3009872, URL https://
doi.org/10.1145/3009837.3009872
[40] Levy, P.B.: Adjunction models for call-by-push-value with stacks. Electronic
Notes in Theoretical Computer Science 69, 248–271 (2003), ISSN 1571-0661,
https://doi.org/https://doi.org/10.1016/S1571-0661(04)80568-1, cTCS’02,
Category Theory and Computer Science
[41] Levy, P.B.: Call-by-push-value: A Functional/Imperative Synthesis, vol. 2.
Springer Netherlands (2003), https://doi.org/10.1007/978-94-007-0954-6
[42] Liang, S., Hudak, P., Jones, M.: Monad transformers and modular inter-
preters. In: ACM SIGPLAN-SIGACT Symposium on Principles of Pro-
gramming Languages, pp. 333–343, POPL ’95, ACM (1995), ISBN 0-89791-
692-1, https://doi.org/10.1145/199448.199528
[43] Lindley, S., McBride, C., McLaughlin, C.: Do be do be do. In:
Proceedings of the 44th ACM SIGPLAN Symposium on Princi-
ples of Programming Languages, p. 500–514, POPL 2017, Associa-
tion for Computing Machinery, New York, NY, USA (2017), ISBN
9781450346603, https://doi.org/10.1145/3009837.3009897, URL https://
doi.org/10.1145/3009837.3009897
[44] Mac Lane, S.: Categories for the Working Mathematician, 2nd edn. Grad-
uate Texts in Mathematics, Springer, Berlin (1998)
[45] Maguire, S.: polysemy: Higher-order, low-boilerplate free monads (2019),
https://hackage.haskell.org/package/polysemy
[46] Matthes, R., Uustalu, T.: Substitution in non-wellfounded syntax with
variable binding. Theoretical Computer Science 327(1), 155–174 (2004),
ISSN 0304-3975, https://doi.org/https://doi.org/10.1016/j.tcs.2004.07.025,
selected Papers of CMCS ’03
[47] Moggi, E.: An abstract view of programming languages. Tech. Rep. ECS-
LFCS-90-113, Edinburgh University, Department of Computer Science
(June 1989)
[48] Moggi, E.: Notions of computation and monads. Informa-
tion and Computation 93(1), 55 – 92 (1991), ISSN 0890-5401,
https://doi.org/https://doi.org/10.1016/0890-5401(91)90052-4, selections
from 1989 IEEE Symposium on Logic in Computer Science
[49] Pir´og, M., Schrijvers, T., Wu, N., Jaskelioﬀ, M.: Syntax and semantics for
operations with scopes. In: Dawar, A., Gr¨adel, E. (eds.) Proceedings of the
33rd Annual ACM/IEEE Symposium on Logic in Computer Science, LICS
2018, Oxford, UK, July 09-12, 2018 (2018)
[50] Pir´og, M., Gibbons, J.: Tracing monadic computations and represent-
ing eﬀects. Electronic Proceedings in Theoretical Computer Science 76,
90–111 (Feb 2012), ISSN 2075-2180, https://doi.org/10.4204/eptcs.76.8,
URL http://dx.doi.org/10.4204/EPTCS.76.8
[51] Plotkin, G., Power, J.: Notions of computation determine monads. In:
Nielsen, M., Engberg, U. (eds.) Foundations of Software Science and
Computation Structures, 5th International Conference, pp. 342–356, FOS-
31
SACS 2002, Springer (2002), https://doi.org/10.1007/3-540-45931-6 24,
URL https://doi.org/10.1007/3-540-45931-6_24
[52] Plotkin, G., Power, J.: Algebraic operations and generic eﬀects. Ap-
plied Categorical Structures 11(1), 69–94 (2003), ISSN 1572-9095,
https://doi.org/10.1023/A:1023064908962
[53] Plotkin, G., Pretnar, M.: A logic for algebraic eﬀects. In: 2008 23rd An-
nual IEEE Symposium on Logic in Computer Science, pp. 118–129 (2008),
https://doi.org/10.1109/LICS.2008.45
[54] Plotkin, G., Pretnar, M.: Handlers of algebraic eﬀects. In: Castagna,
G. (ed.) Programming Languages and Systems, pp. 80–94, Springer
Berlin Heidelberg, Berlin, Heidelberg (2009), ISBN 978-3-642-00590-9,
https://doi.org/10.1007/978-3-642-00590-9 7
[55] Plotkin, G., Pretnar, M.: Handling algebraic eﬀects. Logical Meth-
ods in Computer Science 9(4) (Dec 2013), ISSN 1860-5974,
https://doi.org/10.2168/lmcs-9(4:23)2013
[56] Poulsen, C.B., van der Rest, C., Schrijvers, T.: Staged eﬀects and handlers
for modular languages with abstraction (2021), to Appear
[57] Power, A.J.: Enriched lawvere theories. Theory and Applications of Cate-
gories 6(7), 83–93 (1999)
[58] Riehl, E.: Category Theory in Context. Aurora: Dover Modern Math Orig-
inals, Dover Publications (2017), ISBN 9780486820804
[59] Rix, R., Thomson, P., Wu, N., Schrijvers, T.: fused-eﬀects: A fast, flexi-
ble, fused eﬀect system (2018), https://hackage.haskell.org/package/
fused-effects
[60] Robinson, E.: Variations on algebra: Monadicity and generalisa-
tions of equational theories 13(3), 308–326 (2002), ISSN 1433-299X,
https://doi.org/10.1007/s001650200014
[61] Schrijvers, T., Pir´og, M., Wu, N., Jaskelioﬀ, M.: Monad transformers
and modular algebraic eﬀects: what binds them together. In: Proceed-
ings of the 12th ACM SIGPLAN International Symposium on Haskell,
Haskell@ICFP 2019, Berlin, Germany, August 18-23, 2019, pp. 98–113
(2019), https://doi.org/10.1145/3331545.3342595, URL https://doi.org/
10.1145/3331545.3342595
[62] Spivey, J.M.: Algebras for combinatorial search. Jour-
nal of Functional Programming 19(3-4), 469–487 (2009),
https://doi.org/10.1017/S0956796809007321
[63] Spivey, M.: A functional theory of exceptions. Science of Com-
puter Programming 14(1), 25–42 (1990), ISSN 0167-6423,
https://doi.org/10.1016/0167-6423(90)90056-J
[64] Stark, I.: Free-algebra models for the π-calculus. Theoreti-
cal Computer Science 390(2), 248–270 (2008), ISSN 0304-3975,
https://doi.org/https://doi.org/10.1016/j.tcs.2007.09.024, foundations
of Software Science and Computational Structures
[65] Staton, S.: Instances of computational eﬀects: An algebraic perspective. In:
2013 28th Annual ACM/IEEE Symposium on Logic in Computer Science,
pp. 519–519 (2013), https://doi.org/10.1109/LICS.2013.58
32
[66] Steele, Jr., G.L.: Building interpreters by composing monads. In: Boehm,
H., Lang, B., Yellin, D.M. (eds.) Proceedings of the 21st ACM SIGPLAN-
SIGACT symposium on Principles of programming languages, pp. 472–492,
POPL ’94, ACM (1994), https://doi.org/10.1145/174675.178068
[67] Swierstra, W.: Data types `a la carte. J. Funct. Program. 18(4), 423–436
(2008), https://doi.org/10.1017/S0956796808006758, URL https://doi.
org/10.1017/S0956796808006758
[68] Takano, A., Meijer, E.: Shortcut deforestation in calculational form.
In: Proceedings of the Seventh International Conference on Functional
Programming Languages and Computer Architecture, Association for
Computing Machinery, New York, NY, USA (1995), ISBN 0897917197,
https://doi.org/10.1145/224164.224221, URL https://doi.org/10.1145/
224164.224221
[69] Wadler, P.: Deforestation: Transforming programs to eliminate trees.
Theor. Comput. Sci. 73(2), 231–248 (Jan 1988), ISSN 0304-3975,
https://doi.org/10.1016/0304-3975(90)90147-A, URL https://doi.org/
10.1016/0304-3975(90)90147-A
[70] Wadler, P.: Comprehending monads. In: Proceedings of the 1990 ACM
Conference on LISP and Functional Programming, pp. 61–78, LFP ’90,
ACM (1990), ISBN 0-89791-368-X, https://doi.org/10.1145/91556.91592,
URL https://doi.org/10.1145/91556.91592
[71] Wadler, P.: Monads for functional programming. In: Advanced
Functional Programming, First International Spring School on
Advanced Functional Programming Techniques-Tutorial Text, p.
24–52, Springer-Verlag, Berlin, Heidelberg (1995), ISBN 3540594515,
https://doi.org/10.5555/647698.734146
[72] Wu, N., Schrijvers, T.: Fusion for free. In: Hinze, R., Voigtl¨ander, J. (eds.)
Mathematics of Program Construction, pp. 302–322, Springer International
Publishing, Cham (2015), https://doi.org/978-3-319-19797-5 15
[73] Wu, N., Schrijvers, T., Hinze, R.: Eﬀect handlers in scope. SIGPLAN Not.
(2014)
[74] Yang, Z., Wu, N.: Reasoning about eﬀect interaction by fusion. Proc. ACM
Program. Lang. 5(ICFP) (Aug 2021), https://doi.org/10.1145/3473578,
URL https://doi.org/10.1145/3473578
33
A Figures
↓
KFn
Ix
KEM
Fn
KIx
EM
FreeFn
↑
⊣
⊣
Fn-Alg Endof (C) ×C C
UFn
FreeIx
Ix-Alg C|N| C
UIx
⊣
⊣
FreeP
CP C
⊣
UP
Fig. 2: The resolutions of functorial, indexed and Eilenberg-Moore algebras
data EMAlg Σ Γ x = EM {callEM :: Σ x →x, enterEM :: Γ (Prog Σ Γ x) →x }
handleEM :: (Functor Σ, Functor Γ) ⇒(EMAlg Σ Γ x) →(a →x) →Prog Σ Γ a →x
handleEM alg gen (Return x) = gen x
handleEM alg gen (Call op) = (callEM alg·fmap (handleEM alg gen)) op
handleEM alg gen (Enter op) = (enterEM alg·fmap (fmap (handleEM alg gen))) op
Fig. 3: Haskell implementation of EM algebras of P based on Theorem 3
data Nat= Zero |Nat + 1
data IxAlg Σ Γ a =
IxAlg {action :: ∀n. Σ (a n) →a n
, demote :: ∀n. Γ (a (n + 1)) →a n
, promote :: ∀n. a n →a (n + 1)}
hfold :: (Functor f, Functor g) ⇒IxAlg f g a →∀n :: Nat. Prog f g (a n) →a n
hfold ixAlg (Return x) = x
hfold ixAlg (Call op) = action ixAlg (fmap (hfold ixAlg) op)
hfold ixAlg (Enter scope)
= demote ixAlg (fmap (hfold ixAlg·fmap (promote ixAlg·hfold ixAlg)) scope)
Fig. 4: The hybrid fold for interpreting monad P with indexed algebras [49]
B Notation
This section summarises the mathematical notation used in this paper.
34
Typeface Conventions We use boldface variable such as C and D for abstract
categories and typewriter font for specific categories such as Set. Functors and
objects are denoted by capitalised letters such as F, K, G, X, A in general, but
functors representing signatures of operations are always denoted by Greek let-
ters Γ and Σ, and some concrete functors are denoted by typewriter font such
as Free and Fn. Morphisms and natural transformations are denoted by uncap-
italised letters such as f, g and h and Greek letters such as α, β and τ.
Categories For any two categories C and D, we write C ×D for their product
category, whose objects are denoted by ⟨X, Y⟩where X : C and Y : D and
morphisms are also denoted by ⟨f, g⟩where f : X →X′ and g : Y →Y′. The
functor category from C to D is denoted by DC. Specifically, the category of fini-
tary endofunctors on C is denoted by Endof(C). Given a monad P : C →C, the
category of Eilenberg-Moore algebras is denoted by CP. Given an endofunctor
F on C, the category of F-algebras is F-Alg.
Functors An adjoint pair of functors L : C →D and R : D →C is denoted
by L ⊣R where L is the left adjoint. The associated unit of the adjunction is
η : Id →RL and the counit is ϵ : LR →Id, where Id is the identity functor.
Sometimes we write ηRL or ϵRL to make clear the monad for (co)units.
For monads, letters P, T, M are used, and µis used for their multiplication.
In particular, P is used for the monad of programs with algebraic and scoped
operations, while T is used for the monad from the adjunction FreeFn ↑⊣↓UFn.
Functor G denotes the functor encoding the grammar (13), and functor I
(15) and Fn (16) are used for defining functorial algebras. Similarly, functors Ix,
▹ and ◃ are used for defining indexed algebras (Section 4.2).
In some of the examples in the paper, we make use of the functor List :
Set →Set on Set that sends every set X to the set of finite lists with elements
from X. We use xs ++ys for list concatenation and cons x xs for a list with x in
front of xs.
Comparison functors (Section 4) are denoted by K, and in particular KEM
Fn
and KFn
Ix.
Objects and Morphisms We write f·g for vertical composition of morphisms,
and F◦Gfor horizontal composition of natural transformations and composition
of functors.
Products in a category are denoted by X1 ×···×Xn for finite products
or i∈I Xi for some set I, and coproducts are denoted by X1 +···+ Xn or
i∈I Xi for some set I. The injection morphisms into coproducts are denoted
by ιi : Xi →X1 +···+ Xn for each i, and the morphisms out of coproducts are
[f1,..., fn] : X1 +···+ Xn →Z for objects Z with morphisms fi : Xi →Z for
all i.
Initial objects are generally written as 0, and the unique morphism from the
initial object to an object X is denoted by ! : 0 →X. In particular, the initial
object 0 in a functor category is the constant functor mapping to the initial
object in the codomain.
35
The (carrier of the) initial algebra of an endofunctor G: C →C is denoted
by µG or µY. GY. The structure map of the initial algebra is denoted by in :
G(µG) →G and its inverse is in◦. The unique morphism from the initial G-
algebra to a G-algebra α: GX →X, i.e. the catamorphism for α, is denoted by
α : µG→X.
Given an endofunctor Σ : C →C, the free monad over it is denoted by Σ∗
.
Given a Σ-algebra ⟨X : C, α : ΣX →X⟩and a morphism g : A →X, the
interpretation morphism is denoted by handle⟨X,α⟩g: Σ∗A→X. The notation
handle is also used for functorial and indexed algebras. The structure map of
the free algebra Σ∗A is denoted by op : ΣΣ∗A→Σ∗A.
Functorial algebras are usually denoted by ⟨H, X, αG
, αI⟩ (Definition 1),
where H is an endofunctor on the base category C and X is an object in C.
C Proofs Omitted in the Main Paper
This section contains the proofs omitted in the paper. The following two lemmas
are a more detailed version of the standard results Lemma 1. We refer the reader
to [2, 5] for their proofs.
Lemma 10. If category C has an initial object and colimits of all ω-chains and
endofunctor Σ : C →C preserves them, then there exists an initial Σ-algebra
⟨µΣ : C, in : Σ(µΣ) →µΣ⟩. Moreover, the structure map in is an isomorphism.
Lemma 11. If category C has finite coproducts and colimits of all ω-chains and
endofunctor Σ : C →C preserves them, then the forgetful functor Σ-Alg →C
has a left adjoint FreeΣ : C →Σ-Alg mapping every object X to µY. X+ ΣY
with structure map
opX = Σ(µY. X+ ΣY)ι2 −→X+ Σ(µY. X+ ΣY) in
−→µY. X+ ΣY
and the unit η and counit ϵ of the adjunction FreeΣ ⊣UΣ are
ηX = in·ι1 : X →UΣ(FreeΣX) ϵ⟨X,α⟩= [id, α] : FreeΣX →⟨X, α⟩
(29)
where [id, α] denotes the catamorphism from the initial algebra µY. X+ ΣY
to X. Moreover, this adjunction is strictly monadic.
Lemma (2). The endofunctor Fn (16) on Endof(C) ×C has free algebras, i.e.
there is a functor FreeFn : Endof(C) ×C →Fn-Alg left adjoint to the forgetful
functor UFn : Fn-Alg →Endof(C) ×C.
Proof. Since C is lfp, it is cocomplete, and it follows that Endof(C) ×C is
cocomplete because colimits in functor categories and product categories can be
computed pointwise [44, Thm V.3.2]. It is also easy verification that Fn preserves
all colimits of ω-chains following from the fact that G, Σ and Γ, and all functors
in Endof(C) preserve colimits of ω-chains in their domains respectively. Hence
by Lemma 1, there is a functor FreeFn left adjoint to UFn.
36
Lemma (3). There is a natural isomorphism between FreeFn and the following
functor
FreeFn⟨H, X⟩= G∗H : Endof(C), (IG∗H)∗X : C, opG∗
H , op(IG∗H)∗
X
where opG∗
H : G(G∗H) →G∗H and op(IG∗H)∗
X : IG∗H((IG∗H)∗X) →(IG∗H)∗X
are the structure maps of the free G-algebra and IG∗H-algebra respectively.
Proof. By Lemma 11, the free Fn-algebra generated by ⟨H, X⟩ can be con-
structed from the initial algebra of a functor Fn⟨H,X⟩: Endof(C)×C →Endof(C)×
C such that
Fn⟨H,X⟩Y= ⟨H, X⟩+ FnY
Then we show that an initial Fn⟨H,X⟩-algebra carried by ⟨G∗H, (IG∗H)∗X⟩with
structure map
⟨i1, i2⟩: Fn⟨H,X⟩⟨G∗H, (IG∗H)∗X⟩= ⟨H+ GH, X+ IG∗H((IG∗H)∗X)⟩
→⟨G∗H, (IG∗H)∗X⟩
where i1 = [ηG∗
H , opG∗
H ] : H+ G(G∗H) →G∗H and
i2 = [η(IG∗H)∗
X , op(IG∗H)∗
X ] : X+ IG∗H((IG∗H)∗X) →(IG∗H)∗X.
To see that this Fn⟨H,X⟩-algebra is initial, consider any ⟨C, D⟩in CC ×C with
structure map ⟨j1, j2⟩: Fn⟨H,X⟩⟨C, D⟩→⟨C, D⟩. We have
⟨k1, k2⟩: Fn⟨H,X⟩-Alg ⟨G∗H, (IG∗H)∗X⟩, ⟨i1, i2⟩, ⟨C, D⟩, ⟨j1, j2⟩
⇔ k1 ∈(H+ G−)-Alg ⟨G∗H, i1⟩, ⟨C, j1⟩
∧ k2 ∈(X+ IG∗H−)-Alg ⟨IG∗HX, i2⟩, ⟨D, j2·(X+ Ik1 id)⟩
⇔k1 = ⌈j1·ι1⌉⟨C,j1·ι2 ⟩ ∧ k2 = ⌈j2·(X+ Ik1 )·ι1⌉⟨D,j2·(X+Ik1 )·ι2 ⟩
where we use subscripts of ⌈·⌉to indicate the B for some ⌈f⌉ : LA → B.
The calculation shows that the Fn⟨H,X⟩-algebra homomorphism (k1, k2) uniquely
exists, and thus ⟨G∗H, (IG∗H)∗X⟩with structure map ⟨i1, i2⟩is initial. Then by
Lemma 11, this initial algebra gives the free Fn-algebra generated by ⟨H, X⟩,
and thus we have the isomorphism between FreeFn and FreeFn in the lemma.
This characterisation of free Fn-algebras also allows us to express the unit and
counit of the adjunction FreeFn ⊣UFn in terms of those of some simpler adjunc-
tions.
Lemma 12. Letting the φ be the isomorphism in Lemma 3, the unit of adjunc-
tion FreeFn ⊣UFn is
⟨ηG∗
(IG∗H)∗
H ,η
X ⟩ φ−1η⟨H,X⟩= ⟨H, X⟩ ⟨G∗H, (IG∗H)∗X⟩ Fn∗⟨H, X⟩
37
and its counit ϵ at some Fn-algebra ⟨⟨H, X⟩, ⟨β1, β2⟩⟩is
FreeFn⟨H, X⟩ φ −→ ⟨G∗H, (IG∗H)∗X⟩, ⟨opG∗
H , op(IG∗H)∗
X ⟩ ⟨e1 ,e2 ⟩ −−−−→⟨⟨H, X⟩, ⟨β1, β2⟩⟩
where e1 = UG(ϵG∗
⟨H,β1 ⟩) and e2 = UIG∗H(ϵ(IG∗H)∗
⟨X,β2·Ie1 ⟩).
Proof sketch. It can be calculated from (29) and Lemma 3.
Theorem (1). The monad P is isomorphic to T in the category of monads.
Proof. By (14), P∼
= (Σ+ ΓP)∗= (IP)∗as monads, it is suﬃcient to show that
T is isomorphic to (IP)∗as monads. Recall that P= µG∼
= G∗0 as endofunctors.
Let ψ : G∗0 →P be the isomorphism, and let φ be the isomorphism between
Fn∗and UFnFreeFn by Lemma 3, then for all X : C,
TX= ↓(UFn(FreeFn(↑X))) = ↓(Fn∗⟨0, X⟩) ↓φ
∼
= ↓⟨G∗0, ⟨IG∗0⟩∗X⟩ (30)
(Iψ)∗
= (IG∗0)∗X
∼
= (IP)∗X
Thus T is isomorphic to (IP)∗as endofunctors.
What remains is to show that the isomorphism (30) preserves their units and
multiplications. The unit of T is precisely the unit of the adjunction (FreeFn ↑) ⊣
(↓UFn) composed from the adjunctions ↑⊣↓and FreeFn ⊣UFn. Therefore the unit
of T is ηT
X = ↓(ηFn∗
↑X)·η↓↑
X where η↓↑ and ηFn∗ are the units of the adjunctions
↑⊣↓and FreeFn ⊣UFn respectively. Hence by Lemma 12, we have
ηT
X = ↓(ηFn∗
↑X)·η↓↑
X = ↓(ηFn∗
↑X)·id= ↓(φ−1 ·⟨ηG∗
0 , η(IG∗0 )∗
X ⟩)
= ↓(φ−1)·η(IG∗0 )∗
X = ↓φ·(Iψ
∗)−1
·η(IP)∗
X
which shows that the isomorphism (30) preserves the units of T and (IP)∗
.
Proving the preservation of the multiplications of the two monads is also
direct verification but slightly more involved. By definition, µT
X = ↓UFnϵT
FreeFn ↑X
where ϵT is the counit of the adjunction (FreeFn ↑) ⊣(↓UFn) satisfying
ϵT
FreeFn ↑X = ϵFn∗
FreeFn ↑X·FreeFn(ϵ↓↑
Fn∗↑X) = ϵFn∗
FreeFn ↑X·FreeFn⟨!, id⟩
where ! is the unique G-algebra homomorphism from G∗0 to G∗G∗0. Then
by Lemma 12, we have ϵFn∗
FreeFn ↑X = ⟨e1, e2⟩·φ where e1 = UG(ϵG∗
(G∗0, opG∗
0 )) :
G∗G∗0 → G∗0 and e2 = UIG∗G∗0 (ϵ(IG∗G∗0 )∗
b ) where b is the IG∗G∗0-algebra
⟨(IG∗0)∗X, op(IG∗0 )∗
X·Ie1 ⟩. Hence we have
µT
X = ↓UFn(ϵT
FreeFn ↑X)
= UIG∗G∗0 (ϵ(IG∗G∗0 )∗
b ) ·↓φ·↓UFnFreeFn(!, id)
= UIG∗G∗0 (ϵ(IG∗G∗0 )∗
b ) ·↓φ·↓(φ−1 ·⟨G∗!, (IG∗!)∗⟩·φ)
= UIG∗G∗0 (ϵ(IG∗G∗0 )∗
b )·(IG∗!)∗·↓φ
38
where (IG∗!)∗ is a natural transformation from (IG∗0)∗ to (IG∗G∗0)∗. Then by
base functor fusion [23] i.e. the naturality of the free-forgetful adjunction in the
base functor, we have
µT
X = UIG∗G∗0 (ϵ(IG∗G∗0 )∗
b )·(IG∗!)∗·↓φ= UIG∗0 (ϵ(IG∗0 )∗
b′ ) ·↓φ
where b′ is the IG∗0-algebra with the same carrier as b and with b’s structure
map precomposed with IG∗!:
op(IG∗0 )∗
X·Ie1·IG∗! = op(IG∗0 )∗
X·Ie1·! = op(IG∗0 )∗
X·Iid = op(IG∗0 )∗
X
where the second equality follows from that e1·! : G∗0 →G∗0 is a G-algebra
homomorphism, and thus e1·! = id since G∗0 is an initial G-algebra. Finally we
have
µT
X = UIG∗0 (ϵ(IG∗0 )∗
⟨(IG∗0 )∗X,op
X ⟩) ·↓φ
(IG∗0 )∗
I∗
= UIP(ϵ
P
⟨I∗
PX,op(IP)∗
X ⟩)·((Iψ)∗·↓φ)·T((Iψ)∗·↓φ)
which means exactly that the isomorphism (30) preserves the multiplications of
T and (IP)∗
.
Lemma (7). Functor KEM
Fn is a comparison functor from the Eilenberg-Moore
resolution of P to the resolution FreeFn ↑⊣↓UFn of functorial algebras.
Proof. By Definition 2, we need to show that KEM
Fn commutes with left and right
adjoints of both resolutions: for right adjoints, we have UΣ+Γ◦P⟨X, α⟩= X=
↓UFnKEM
Fn ⟨X, α⟩; and for left adjoints,
KEM
Fn (FreeΣ+Γ◦PX) = KEM
Fn ⟨(Σ+ Γ◦P)∗X, op(Σ+Γ◦P)∗
X ⟩
= ⟨P, (Σ+ Γ◦P)∗X, in, op···
X⟩
{P∼
= G∗0 and (Σ+ Γ◦P) = IP (Section 3.3)}
∼
= ⟨G∗0, (IG∗0)∗X, opG∗
0 , op(IG∗0 )∗
X ⟩ {By Lemma 3}
∼
= FreeFn(↑X)
and similarly for the actions on morphisms. Here we only have KEM
Fn FreeΣ+Γ◦P
being isomorphic to FreeFn ↑instead of a strict equality, since these two reso-
lutions induce the monad P only up to isomorphism. To remedy this, one can
slightly generalise the definition of comparison functors to take an isomorphism
into account, but we leave it out here.
Lemma (8). KFn
Ix is a comparison functor from the resolution FreeFn ↑⊣↓UFn
of functorial algebras to the resolution FreeIx ⊣ UIx of indexed algebras.
Proof. We need to show the required commutativities for KFn
Ix to be a comparison
functor:
↓UFn
∼
= UIxKFn
Ix and KFn
Ix FreeFn ↑∼
= FreeIx
39
First it is easy to see that it commutes with the right adjoints:
UIx(KFn
Ix ⟨H, X, αG
, αI⟩) = UIx⟨A, a, d, p⟩= A0 = X= UIx(KFn
Ix ⟨H, X, αG
, αI⟩)
Its commutativity with the left adjoints is slightly more involved, and we show
a sketch here. Pir´og et al. [49] show that FreeIx X is isomorphic to the indexed
algebras carried by P+
X : C|N|such that (P+
X)i = Pi+1X with structure maps
kΓ▹
n
kΣ
n
= Σ(P+
X)n = ΣPPnX in·ι2 −−−→PPnX
= (Γ ▹P+
X)n = ΓPPPnX in·ι3 −−−→PPnX
k◃
n
= (P+
X)n = PPnX in·ι1 −−−→PPPnX
where in : Id + Σ◦P + Γ◦P◦P → P is the isomorphism between P and
GP. Also by Lemma 3, we know that FreeFn ↑X is isomorphic to the functorial
algebra ⟨P, PX, in, [in·ι2, in·ι3]⟩. Clearly KFn
Ix FreeFn ↑and FreeIx agree on
the carrier P+
X. It can be checked that they agree on the structure maps and the
action on morphisms too.
D Handling with EM Algebras in Haskell
This section provides some details about how Eilenberg-Moore algebras of the
monad P (11) are used for interpret scoped operations. First we have the follow
theorem characterising handle in a way that suits recursive implementation. A
Haskell implementation is in (Figure 3).
Theorem 3 (Interpreting with EM Algebras). Given an Eilenberg-Moore
algebra as in (21), for any morphism g: A→X in C for some A, the interpre-
tation of PA with this algebra and g satisfies
handle⟨X,αΣ,αΓ⟩g= [ g, αΣ·Σ(handle⟨X,αΣ,αΓ⟩g), (31)
αΓ·ΓP(handle⟨X,αΣ,αΓ⟩g)]·in◦
A
where in◦
: P →Id + Σ◦P + Γ◦P◦P is the isomorphism between P and GP.
Proof. It can be calculated by plugging the counit ϵfor the monad FreeΣ+Γ◦P ⊣
UΣ+Γ◦P (Lemma 11) into (22).
Example 13. With the implementation of EM algebras and their interpretation
(Figure 3), the scoped operation once in Example 4 can be interpreted by the
following EM algebra
onceAlgEM :: EMAlg Choice Once [a ]
onceAlgEM= EM {..}where
callEM :: Choice [a ] →[a ]
callEM Fail = [ ]
callEM (Or x y) = x ++y
40
enterEM :: Once (Prog Choice Once [a ]) →[a ]
enterEM (Once p) =
case handleEM onceAlgEM (λx →[x ]) p of
[ ] →[ ]
(x: ) →x
Note that this EM algebra is defined with recursion: the part of syntax enclosed
by the scope of once is interpreted by a recursive call to handleEM with the
algebra itself.
E Handling with Indexed Algebras by Hybrid Fold
This section elaborates Example 12 about using fusion laws and functorial al-
gebra to prove the correctness of a recursive scheme used by Pir´og et al. [49] to
interpret scoped operations with indexed algebras. Although Pir´og et al. [49] pro-
pose the adjunction (FreeIx ) ⊣( UIx) for interpreting scoped operations with
indexed algebras, they use the recursive function hfold (Figure 4) in their im-
plementation to interpret P (12) with indexed algebras. Compared to a faithful
implementation of handle, their hfold is more eﬃcient since it skips transforming
P to the free indexed algebra ( UIxFreeIx ) but directly works on P. Thus we
call it a hybrid fold since it works on a syntactic structure P that is not freely
generated by its type of algebras giving semantics.
While the definition of hfold is computationally intuitive, Pir´og et al. [49] did
not provide formal justification for this recursive definition. In this subsection, we
fill the gap by showing that hfold coincides with handle with indexed algebras.
We divide the proof into three parts for clarity: after making the problem precise
(Section E.1), we first show that the hfold for an indexed algebra A is equivalently
a catamorphism from P in Endof(C) (Section E.2), which is then a special case
of interpreting with functorial algebras (Section E.3), and finally we translate
this functorial algebra into the category Ix-Alg of indexed algebras using KFn
Ix ,
and show that it induces the same interpretation as the one from the indexed
algebra A that we start with (Section E.4).
E.1 Semantic Problem of Hybrid Fold
Fix an indexed algebra carried by A: C|N|in this section:
⟨A: C|N|
, a: Σ◦A→A, d: Γ◦(▹A) →A, p: A→(▹A)⟩ (32)
For notational convenience, we define a functor S : |N|→|N|such that Sn=
n+ 1, then ▹A= A◦S since (▹A)n = An+1 = A(Sn). With functor S, we
can view p and d as p : A →A◦S and d : Γ◦A◦S →A. Then the recursive
definition of hfold in Figure 4 can be understood as a morphism h: P◦A→A
in C|N|satisfying the equation
h= [h1, h2, h3]·(in◦◦A) (33)
41
where in◦
: P →Id + Σ◦P+ Γ◦P◦P is the isomorphism between P and GP,
and h1, h2 and h3 correspond to the three cases of hfold respectively:
h1 = Id ◦A id
−→A h2 = Σ◦P◦A Σ◦h
−−−→Σ◦A a
−→A
h3 = Γ◦P◦P◦A Γ◦P◦h
−−−−−→Γ◦P◦A Γ◦P◦p −−−−−→Γ◦P◦A◦S Γ◦h◦S
−−−−−→Γ◦A◦S d
−→A
In general, an equation in the form of (33) does not necessarily have a (unique)
solution. Thus the semantics of the function hfold is not clear. We settle this
problem with the following result.
Theorem 4 (Hybrid Folds Coincide with Interpretation). There exists
a unique solution to (33) and it coincides with the interpretation with indexed
algebra A at level 0 (26):
h0 = handle⟨A,a,d,p⟩id : PA0 →A0
We prove the theorem in the rest of this section with the tools that we have
developed.
E.2 Hybrid Fold Is an Adjoint Fold
The first step of our proof is to show the unique existence of the solution to
(33) based on the observation that it is an adjoint fold equation [23] with the
adjunction between right Kan extension [44] and composition with A. Hinze [23]
shows the following theorem stating that adjoint fold equations have a unique
solution.
Theorem 5 (Mendler-style Adjoint Folds [23]). Given any adjunction
L⊣R: D →C, an endofunctor G: D →D whose initial algebra ⟨µG, in⟩exists,
and a natural transformation Φ: C(L−
, B) →C(LD−
, B) for some B : C, then
there exists a unique x: L(µG) →B satisfying
x= ΦµG(x)·Lin◦ (34)
and the unique solution satisfies ⌊x⌋= ⌊ΦRB(ϵB)⌋ where ⌊·⌋: C(LD, C) →
D(D, RC) is the isomorphism for the adjunction L⊣R.
Since h : PA →A and P= µG : Endof(C), to apply this theorem to (33),
we only need to (i) make (−◦A) : Endof(C) →C|N|a left adjoint and (ii) make
[h1, h2, h3] in (33) an instance of ΦP(h) for some natural transformation Φ:
– For (i), the functor−◦Ais left adjoint to the right Kan extension along A,
that is a functor RanA : C|N|→Endof(C), which always exists when C is lfp.
– For (ii), we define a natural transformation Φ: C|N|(−◦A, A) →C|N|((G−)◦A, A)
such that for all H : Endof(C) and f ∈C|N|(H◦A, A),
ΦH(f) = [f1, f2, f3]·(in◦◦A)
42
where
f1 = Id ◦A id
−→A f2 = Σ◦H◦A Σ◦f −−−→Σ◦A a
−→A (35)
f3 = Γ◦H◦H◦A Γ◦H◦f −−−−−→Γ◦H◦A Γ◦H◦p −−−−−→Γ◦H◦A◦S Γf −−→Γ◦A◦S d
−→A
It is immediate that hfold (33) is exactly h= ΦP(h)·(in◦◦A).
Then by Theorem 5 we have the following result.
Lemma 13 (Unique Existence of Hybrid Fold). The recursive definition
h: PA→A (33) of hybrid folds (Figure 4) has a unique solution:
h= PA ⌊ΦRanAA(ϵA)⌋A
−−−−−−−−−−−→(RanAA)A ϵA −→A (36)
where ϵA : (RanAA)◦A →A is the counit, and ⌊ΦRanAA(ϵA)⌋ is the catamor-
phism from the initial G-algebra P (12) to the G-algebra carried by RanAA with
structure map ⌊ΦRanAA(ϵA)⌋.
E.3 Catamorphism as Interpretation
We have shown that hfoldA for all indexed algebras A uniquely exists, and now
to show its coincidence with handleA (the second part of Theorem 4), we show
that the hybrid fold coincides with handle with some functorial algebra, and
then use the comparison functor KFn
Ix (Section 4.5) to translate this functorial
algebra into an indexed algebra.
The following lemma relating G-Alg and Fn-Alg is straightforward.
Lemma 14. For the functor G in (13) and every X : C, there is a faithful
functor JX : G-Alg →Fn-Alg that maps any G-algebra ⟨H, α: GH →H⟩to the
functorial algebra ⟨H, HX, α, β⟩where β = [αX·ι2, αX·ι3] : IH(HX) →HX.
Now consider the G-algebra in Lemma 13 carried by RanAA with structure
map αG = ⌊ΦRanAA(ϵA)⌋, we can trivially make it a functorial algebra ˆ α by
Lemma 14, hoping for applications of the fusion law of functorial algebras (5.1):
ˆ
α= JA0 ⟨RanAA, αG⟩= ⟨RanAA, (RanAA)A0, αG
, αG
·[ι2, ι3]⟩ (37)
The catamorphism to αG and handleˆ
α are related by αG A0
= handleˆ
α (αG
A0
·ι1),
which can be shown by checking that the formula (20) for computing handleˆ
α is
exactly the defining equation of the catamorphism αG A0 . Plugging the identity
into (36), we obtain
h0 = (ϵA)0·handleˆ
α (αG
A0
·ι1) (38)
Now we have made some progress because the right-hand side takes exactly the
form for applications of the fusion law—an interpretation followed by a morphism
(ϵA)0 : (RanAA)A0 →A0.
In order to apply the fusion law, we need to make (ϵA)0 : (RanAA)A0 →A0
a functorial algebra homomorphism from ˆ α(37). With some exploration, we can
43
ˆ
find a functorial algebra
α′ = ⟨RanAA, A0, αG
, αI⟩diﬀering from ˆ α by the
second and fourth components, where
αI = [a0, d0·Γ(ϵA)0·Γ(RanAA)p] : ΣA0 + Γ((RanAA)A0) →A0
and a, d and p are the structure maps of indexed algebra A (32). It can be
checked that ⟨id, (ϵA)0⟩is a functorial algebra homomorphism from ˆ α to
ˆ
α′
.
Thus by Equation 5.1, we obtain
h0 = handleˆ
α (ϵA)0·(αG)A0·ι1 = handleˆ
α′ id : PA0 →A0 (39)
which means that the hybrid fold coincides with the interpretation with functo-
rial algebra
ˆ
α′
.
E.4 Translating Back to Indexed Algebras
ˆ
The last step of our proof is translating the functorial algebra
α′ back to an
indexed algebra using comparison functor KFn
Ix (Section 4.5), and showing that
the resulting indexed algebra induces the same interpretation morphism as the
one induced by A.
Recall that the comparison functor KFn
Ix maps a functorial algebra carried by
ˆ
⟨H, X⟩to an indexed algebra carried by i →HiX. Thus KFn
α′ is an indexed
Ix
algebra carried by i→(RanAA)iA0 and by Lemma 5 and (39), KFn
Ix preserves the
induced interpretation:
h0 = handleˆ
α′ id= handleKFn
ˆ
α′ id : PA0 →A0 Ix
(40)
What remains is to prove that handleKFn
ˆ
α′ id= handleA id, and we show this
Ix
by the fusion law (of indexed algebras): define a natural transformation τ :
(RanAA)iA0 →Ai between C|N|functors by τ0 = id : (RanAA)0A0 →A0 and
τi+1 = (RanAA)(RanAA)iA0
(RanAA)τi
−−−−−−→(RanAA)Ai
(RanAA)p
−−−−−−→(RanAA)Ai+1
ϵA −→Ai+1
ˆ
and it can be checked that τ is an indexed algebra homomorphism from KFn
α′
Ix
to ⟨A, a, d, p⟩. Note that we have UIxτ equals id : A0 →A0, and by Lemma 5,
we have that
h0 = handleKFn
ˆ
α′ id= id·handleKFn
ˆ
α′ id = ( UIx)τ·handleKFn
ˆ
α′ id= handleA id
Ix
Ix
Ix
This completes our proof of Theorem 4 saying that Pir´og et al. [49]’s hfold indeed
correctly implements handle with indexed algebras.
To summarise, we first connected hfold to a functorial algebra using a Mendler-
style adjoint fold, and then used the fusion law of functorial algebras to simplify
it. Then it is translated to the category of indexed algebras, and we used the
fusion law one more time there.